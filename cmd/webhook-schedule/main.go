package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/silasstoffel/account-service/configs"
	"github.com/silasstoffel/account-service/internal/domain/webhook"
	"github.com/silasstoffel/account-service/internal/event"
	"github.com/silasstoffel/account-service/internal/infra/database"
	"github.com/silasstoffel/account-service/internal/infra/helper"
	"github.com/silasstoffel/account-service/internal/infra/messaging"
	"github.com/silasstoffel/account-service/internal/logger"
)

var message event.Event

type sqsSender struct {
	sqsClient *sqs.Client
	queueUrl  string
}

type scheduleMessageInput struct {
	MessageId      string `json:"messageId"`
	EventId        string `json:"eventId"`
	SubscriptionId string `json:"subscriptionId"`
	EventType      string `json:"eventType"`
	Data           string `json:"data"`
}

func main() {
	log.Println("Starting webhook schedule consumer")
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
	consumer := messaging.MessagingConsumer{
		SqsClient:           snsClient,
		QueueUrl:            config.Aws.WebhookScheduleQueueUrl,
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     3,
		VisibilityTimeout:   20,
	}

	scheduleSenderConfig := sqs.NewFromConfig(awsConfig)
	scheduleSender := sqsSender{
		sqsClient: scheduleSenderConfig,
		queueUrl:  config.Aws.WebhookSenderQueueUrl,
	}
	logger := logger.NewLogger(config)
	subscriptionRepository := database.NewSubscriptionRepository(cnx, logger)
	messageChannel := make(chan *types.Message)

	go consumer.PollingMessages(messageChannel)

	for rawMessage := range messageChannel {
		fmt.Println("Processing message", *rawMessage.MessageId)
		err := messaging.ExtractMessageFromTopic(rawMessage, &message)
		if err != nil {
			fmt.Println("Error parsing message", err)
			continue
		}

		var event event.Event
		if err := dataMessageToEvent(&message.Data, &event); err != nil {
			fmt.Println("Error when convert message to event", err)
			continue
		}

		subscriptions, err := subscriptionRepository.GetByEventType(event.Type)
		if err != nil {
			fmt.Println("Error when get subscriptions to schedule", err)
			continue
		}

		messageBatch, err := buildMessageBatch(subscriptions, event)
		if err != nil {
			fmt.Println("Error when build message batch", err)
			continue
		}
		err = scheduleSender.schedule(messageBatch)
		if err != nil {
			fmt.Println("Error when schedule messaging", err)
			continue
		}
		consumer.DeleteMessage(*rawMessage.ReceiptHandle)
		fmt.Println("Processed message:", *rawMessage.MessageId)
	}
}

func buildMessageBatch(subscriptions []webhook.Subscription, event event.Event) ([]types.SendMessageBatchRequestEntry, error) {
	var messageBatch []types.SendMessageBatchRequestEntry
	counter := 1
	for _, subscription := range subscriptions {
		messageId := helper.NewULID()
		data := &scheduleMessageInput{
			MessageId:      messageId,
			EventId:        event.Id,
			SubscriptionId: subscription.Id,
			EventType:      event.Type,
			Data:           event.Data,
		}
		messageBody, err := json.Marshal(data)

		if err != nil {
			detail := "Error when convert event payload to json."
			log.Println(detail, "Detail", err.Error())
			return messageBatch, err
		}

		messageBatch = append(messageBatch, types.SendMessageBatchRequestEntry{
			Id:          aws.String(strconv.Itoa(counter)),
			MessageBody: aws.String(string(messageBody)),
		})
		counter++
	}
	return messageBatch, nil
}

func (ref *sqsSender) schedule(entries []types.SendMessageBatchRequestEntry) error {
	if len(entries) == 0 {
		return nil
	}
	_, err := ref.sqsClient.SendMessageBatch(context.TODO(), &sqs.SendMessageBatchInput{
		Entries:  entries,
		QueueUrl: aws.String(ref.queueUrl),
	})

	if err != nil {
		fmt.Println("Error sending message", err)
		return err
	}

	return nil
}

func dataMessageToEvent(message *string, event *event.Event) error {
	if err := json.Unmarshal([]byte(*message), event); err != nil {
		message := "Error when convert event payload to json."
		log.Println(message, "Detail", err.Error())
		return err
	}
	return nil
}
