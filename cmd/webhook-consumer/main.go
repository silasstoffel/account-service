package main

import (
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

func main() {
	log.Println("Starting webhook consumer")
	config := configs.NewConfigFromEnvVars()
	awsConfig, err := helper.BuildAwsConfig(config.Aws.Endpoint)
	if err != nil {
		log.Println("Error creating aws config", err)
		panic(err)
	}

	snsClient := sqs.NewFromConfig(awsConfig)
	consumer := messaging.MessagingConsumer{
		SqsClient:           snsClient,
		QueueUrl:            config.Aws.WebhookSenderQueueUrl,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     1,
	}

	messageChannel := make(chan *types.Message, 2)

	go consumer.PollingMessages(messageChannel)

	for rawMessage := range messageChannel {
		fmt.Println("Processing message", rawMessage.MessageId)
		err := messaging.ExtractMessageFromTopic(rawMessage, &message)
		if err != nil {
			fmt.Println("Error parsing message", err)
			continue
		}
		consumer.DeleteMessage(*rawMessage.ReceiptHandle)
		fmt.Println(message.Data)
		fmt.Println("Processed message:", *rawMessage.MessageId)
	}
}
