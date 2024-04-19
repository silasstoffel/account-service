package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/silasstoffel/account-service/configs"
	"github.com/silasstoffel/account-service/internal/event"
	"github.com/silasstoffel/account-service/internal/infra/database"
	"github.com/silasstoffel/account-service/internal/infra/helper"
	"github.com/silasstoffel/account-service/internal/infra/messaging"
	"github.com/silasstoffel/account-service/internal/logger"
	usecase "github.com/silasstoffel/account-service/internal/usecase/event"
)

var message event.Event

func main() {
	log.Println("Starting events consumer")
	config := configs.NewConfigFromEnvVars()
	awsConfig, err := helper.BuildAwsConfig(config)
	if err != nil {
		log.Println("Error creating aws config", err)
		panic(err)
	}

	snsClient := sqs.NewFromConfig(awsConfig)
	consumer := messaging.MessagingConsumer{
		SqsClient:           snsClient,
		QueueUrl:            config.Aws.AccountServiceQueueUrl,
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     1,
		VisibilityTimeout:   30,
	}

	cnx, err := database.OpenConnection(config)
	if err != nil {
		log.Fatalf("Failed to open connection to database: %v", err)
		return
	}
	defer cnx.Close()

	logger := logger.NewLogger(config)
	messagingProducer := messaging.NewDefaultMessagingProducerFromConfig(config)
	eventRepository := database.NewEventRepository(cnx, logger)

	createEventUseCase := usecase.CreateEventUseCase{
		EventRepository: eventRepository,
		Messaging:       messagingProducer,
	}
	messageChannel := make(chan *types.Message)

	go consumer.PollingMessages(messageChannel)

	for rawMessage := range messageChannel {
		err = messaging.ExtractMessageFromTopic(rawMessage, &message)
		if err != nil {
			fmt.Println("Error parsing or extract message", err)
			continue
		}

		err = createEventUseCase.CreateEventUseCase(message)
		if err != nil {
			log.Println("Error creating event", err)
			continue
		}
		consumer.DeleteMessage(*rawMessage.ReceiptHandle)
		fmt.Println("Processed message", *rawMessage.MessageId)
	}
}
