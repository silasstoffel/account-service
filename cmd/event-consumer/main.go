package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/silasstoffel/account-service/internal/event"
	"github.com/silasstoffel/account-service/internal/infra/database"
	"github.com/silasstoffel/account-service/internal/infra/helper"
	"github.com/silasstoffel/account-service/internal/infra/messaging"
	usecase "github.com/silasstoffel/account-service/internal/usecase/event"
)

type MessageSchema struct {
	Message           string
	MessageId         string
	MessageAttributes interface{}
}

func main() {
	log.Println("Starting events consumer")

	awsConfig, err := helper.BuildAwsConfig("http://localhost:4566")
	if err != nil {
		log.Println("Error creating aws config", err)
		panic(err)
	}

	var message event.Event
	var schema MessageSchema

	snsClient := sqs.NewFromConfig(awsConfig)
	consumer := messaging.MessagingConsumer{
		SqsClient:           snsClient,
		QueueUrl:            "http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/account-service",
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     1,
	}

	cnx := database.OpenConnection()
	defer cnx.Close()

	messagingProducer := messaging.NewMessagingProducer(
		"arn:aws:sns:us-east-1:000000000000:account-service-topic",
		"http://localhost:4566",
	)
	eventRepository := database.NewEventRepository(cnx)

	createEventUseCase := usecase.CreateEventUseCase{
		EventRepository: eventRepository,
		Messaging:       messagingProducer,
	}
	messageChannel := make(chan *types.Message, 2)

	go consumer.PollingMessages(messageChannel)

	for rawMessage := range messageChannel {

		if err = json.Unmarshal([]byte(*rawMessage.Body), &schema); err != nil {
			log.Println("Error unmarshalling body", err)
			panic(err)
		}

		if err = json.Unmarshal([]byte(schema.Message), &message); err != nil {
			log.Println("Error unmarshalling message", err)
			panic(err)
		}

		err := createEventUseCase.CreateEventUseCase(message)
		if err != nil {
			log.Println("Error creating event", err)
			continue
		}
		consumer.DeleteMessage(*rawMessage.ReceiptHandle)

		fmt.Println("Processed message", schema.MessageId)
	}
}
