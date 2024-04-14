package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/silasstoffel/account-service/configs"
	"github.com/silasstoffel/account-service/internal/domain/webhook"
	"github.com/silasstoffel/account-service/internal/exception"
	"github.com/silasstoffel/account-service/internal/infra/database"
	"github.com/silasstoffel/account-service/internal/infra/helper"
	"github.com/silasstoffel/account-service/internal/infra/messaging"
)

type subscriptionMessageDetail struct {
	Id        string `json:"id"`
	EventType string `json:"eventType"`
	Url       string `json:"url"`
}

type messageDetail struct {
	MessageId      string                    `json:"messageId"`
	EventId        string                    `json:"eventId"`
	SubscriptionId string                    `json:"subscriptionId"`
	EventType      string                    `json:"eventType"`
	Data           string                    `json:"data"`
	Subscription   subscriptionMessageDetail `json:"subscription"`
	SendAt         time.Time                 `json:"sendAt"`
}

type notifyStats struct {
	statusCode int
	startedAt  time.Time
	finishedAt time.Time
}

var transactionRepository *database.WebhookTransactionRepository

func main() {
	config := configs.NewConfigFromEnvVars()
	awsConfig, err := helper.BuildAwsConfig(config)
	if err != nil {
		log.Println("Error creating aws config", err)
		panic(err)
	}

	cnx, err := database.OpenConnection(config)
	if err != nil {
		log.Fatalf("Failed to open connection to database: %v", err)
		return
	}
	defer cnx.Close()

	snsClient := sqs.NewFromConfig(awsConfig)
	consumer := messaging.NewMessagingConsumer(config.Aws.WebhookSenderQueueUrl, snsClient)
	consumer.VisibilityTimeout = 45
	consumer.WaitTimeSeconds = 10

	transactionRepository = database.NewWebhookTransactionRepository(cnx)

	messageChannel := make(chan *types.Message)

	go consumer.PollingMessages(messageChannel)

	var message messageDetail
	ttl := 3 * time.Second
	for rawMessage := range messageChannel {
		log.Println("Processing message", *rawMessage.MessageId)

		err := messaging.ExtractMessageFromQueue(rawMessage, &message)
		if err != nil {
			log.Println("Error parsing message", err)
			continue
		}

		delete := true
		stats, err := notify(message, ttl)
		if err != nil {
			detail := err.(*exception.Exception)
			if detail.Code == webhook.WebhookTransactionNotificationTimeout {
				delete = false
			}
		}

		err = upsert(message, stats)
		if err != nil {
			log.Println("Error upserting transaction", err)
			continue
		}
		if delete {
			consumer.DeleteMessage(*rawMessage.ReceiptHandle)
		}
		log.Println("Processed message:", *rawMessage.MessageId)
	}
}

func notify(message messageDetail, ttl time.Duration) (notifyStats, error) {
	message.SendAt = time.Now().UTC()
	stats := notifyStats{startedAt: message.SendAt, statusCode: 0}
	payload, err := json.Marshal(message)
	if err != nil {
		stats.finishedAt = time.Now().UTC()
		log.Println("Error marshalling message on notify webhook", err)
		return stats, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		message.Subscription.Url,
		strings.NewReader(string(payload)),
	)
	if err != nil {
		stats.finishedAt = time.Now().UTC()
		log.Println("Error creating request to notify webhook", err)
		return stats, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Message-Id", message.MessageId)
	req.Header.Set("X-Send-At", time.Now().UTC().Format(time.RFC3339))
	req.Header.Set("User-Agent", "account-service/webhook")

	client := &http.Client{Timeout: ttl}
	resp, err := client.Do(req)
	if err != nil {
		stats.finishedAt = time.Now().UTC()
		isTimeout := strings.Contains(err.Error(), "Client.Timeout")
		if isTimeout {
			message := "Webhook notification timeout"
			log.Println(message, err.Error())
			return stats, exception.New(exception.WebhookTransactionNotificationTimeout, &err)
		}
		return stats, err
	}
	defer resp.Body.Close()
	stats.finishedAt = time.Now().UTC()
	stats.statusCode = resp.StatusCode

	return stats, nil
}

func upsert(message messageDetail, stats notifyStats) error {
	_, err := transactionRepository.FindById(message.MessageId)
	createTransaction := false

	if err != nil {
		detail := err.(*exception.Exception)
		if detail.Code != webhook.WebhookTransactionNotFound {
			log.Println("Error finding transaction", err)
			return err
		}
		createTransaction = true
	}

	if createTransaction {
		transaction := webhook.WebhookTransaction{
			Id:                 message.MessageId,
			EventId:            message.EventId,
			SubscriptionId:     message.SubscriptionId,
			EventType:          message.EventType,
			ReceivedStatusCode: stats.statusCode,
			RequestStartedAt:   stats.startedAt,
			RequestFinishedAt:  stats.finishedAt,
			NumberOfRequests:   1,
		}
		_, err := transactionRepository.Create(transaction)
		if err != nil {
			log.Println("Error when creating transaction", err)
			return err
		}
		return nil
	}

	_, err = transactionRepository.Update(message.MessageId, webhook.UpdateTransactionInput{
		ReceivedStatusCode: stats.statusCode,
		RequestStartedAt:   stats.startedAt,
		RequestFinishedAt:  stats.finishedAt,
	})

	if err != nil {
		log.Println("Error when updating transaction", err)
		return err
	}

	return nil
}
