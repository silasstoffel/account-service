package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/silasstoffel/account-service/configs"
	"github.com/silasstoffel/account-service/internal/event"
	"github.com/silasstoffel/account-service/internal/infra/helper"
	"github.com/silasstoffel/account-service/internal/infra/messaging"
)

var message event.Event

type sqsSender struct {
	sqsClient *sqs.Client
	queueUrl  string
}

func main() {
	log.Println("Starting webhook schedule consumer")
	config := configs.NewConfigFromEnvVars()
	awsConfig, err := helper.BuildAwsConfig(config.Aws.Endpoint)
	if err != nil {
		log.Println("Error creating aws config", err)
		panic(err)
	}

	snsClient := sqs.NewFromConfig(awsConfig)
	consumer := messaging.MessagingConsumer{
		SqsClient:           snsClient,
		QueueUrl:            config.Aws.WebhookScheduleQueueUrl,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     1,
	}

	scheduleSenderConfig := sqs.NewFromConfig(awsConfig)
	scheduleSender := sqsSender{
		sqsClient: scheduleSenderConfig,
		queueUrl:  config.Aws.WebhookSenderQueueUrl,
	}

	messageChannel := make(chan *types.Message, 2)

	go consumer.PollingMessages(messageChannel)

	for rawMessage := range messageChannel {
		fmt.Println("Processing message", *rawMessage.MessageId)
		err := messaging.ExtractMessageFromTopic(rawMessage, &message)
		if err != nil {
			fmt.Println("Error parsing message", err)
			continue
		}
		err = scheduleSender.schedule(message)
		if err != nil {
			fmt.Println("Error when schedule messaging", err)
			continue
		}
		consumer.DeleteMessage(*rawMessage.ReceiptHandle)
		fmt.Println("Processed message:", *rawMessage.MessageId)
	}
}

func (ref *sqsSender) schedule(message interface{}) error {
	fmt.Println("Scheduling message")
	dataAsJson, err := json.Marshal(message)
	if err != nil {
		message := "Error when convert event payload to json."
		log.Println(message, "Detail", err.Error())
		return err
	}

	msg := string(dataAsJson)
	output, err := ref.sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: &msg,
		QueueUrl:    &ref.queueUrl,
	})

	if err != nil {
		fmt.Println("Error sending message", err)
		return err
	}
	fmt.Println("Message schedule", *output.MessageId)
	return nil
}
