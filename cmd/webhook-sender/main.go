package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/silasstoffel/account-service/configs"
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
	EventId        string                    `json:"eventId"`
	SubscriptionId string                    `json:"subscriptionId"`
	EventType      string                    `json:"eventType"`
	Data           string                    `json:"data"`
	Subscription   subscriptionMessageDetail `json:"subscription"`
}

type responseStats struct {
	statusCode int
	startedAt  time.Time
	finishedAt time.Time
}

func main() {
	log.Println("Starting webhook sender consumer")
	config := configs.NewConfigFromEnvVars()
	awsConfig, err := helper.BuildAwsConfig(config.Aws.Endpoint)
	if err != nil {
		log.Println("Error creating aws config", err)
		panic(err)
	}

	cnx := database.OpenConnection(config)
	defer cnx.Close()

	snsClient := sqs.NewFromConfig(awsConfig)
	consumer := messaging.NewMessagingConsumer(config.Aws.WebhookSenderQueueUrl, snsClient)
	consumer.VisibilityTimeout = 15
	consumer.WaitTimeSeconds = 10

	messageChannel := make(chan *types.Message)

	go consumer.PollingMessages(messageChannel)

	var message messageDetail
	ttl := 5 * time.Second
	for rawMessage := range messageChannel {
		fmt.Println("Processing message", *rawMessage.MessageId)

		err := messaging.ExtractMessageFromQueue(rawMessage, &message)
		if err != nil {
			fmt.Println("Error parsing message", err)
			continue
		}

		stats, err := notify(message, ttl)
		if err != nil {
			fmt.Println("Error notifying webhook", err)
			continue
		}
		consumer.DeleteMessage(*rawMessage.ReceiptHandle)
		log.Println("Webhook response", stats)
		log.Println("Processed message:", *rawMessage.MessageId)
	}
}

func notify(message messageDetail, ttl time.Duration) (responseStats, error) {
	fmt.Println("Notifying webhook", message.Subscription.Url)
	stats := responseStats{startedAt: time.Now(), statusCode: 0}
	req, err := http.NewRequest(
		http.MethodPost,
		message.Subscription.Url,
		strings.NewReader(message.Data),
	)
	if err != nil {
		return stats, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Event-Id", message.EventId)
	req.Header.Set("X-Event-Type", message.EventType)
	req.Header.Set("X-Subscription-Id", message.SubscriptionId)
	req.Header.Set("X-Source", "account-service")
	req.Header.Set("User-Agent", "account-service/webhook")

	client := &http.Client{
		Timeout: ttl,
	}
	resp, err := client.Do(req)
	stats.finishedAt = time.Now()
	stats.statusCode = resp.StatusCode

	if err != nil {
		log.Println("Error notifying webhook", err)
		return stats, err
	}
	defer resp.Body.Close()

	log.Println("Notified webhook", message.Subscription.Url)
	return stats, nil
}
