package main

import (
	"fmt"

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
	config := configs.NewConfigFromEnvVars()
	logger := logger.NewLogger(config)
	logger.Info("Starting events consumer", nil)
	awsConfig, err := helper.BuildAwsConfig(config)
	if err != nil {
		logger.Error("Error creating aws config", err, nil)
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
		logger.Error("Failed to open connection to database", err, nil)
		panic(err)
	}
	defer cnx.Close()

	messagingProducer := messaging.NewDefaultMessagingProducerFromConfig(config, logger)
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
			logger.Error("Error parsing or extract message", err, nil)
			continue
		}

		err = createEventUseCase.CreateEventUseCase(message)
		if err != nil {
			logger.Error("Error creating event", err, nil)
			continue
		}
		consumer.DeleteMessage(*rawMessage.ReceiptHandle)
		fmt.Println("Processed message", *rawMessage.MessageId)
		logger.Info("Processed message", map[string]interface{}{
			"messageId": *rawMessage.MessageId,
		})
	}
}
