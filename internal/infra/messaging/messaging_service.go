package messaging

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqsType "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/silasstoffel/account-service/internal/event"
	"github.com/silasstoffel/account-service/internal/exception"
	"github.com/silasstoffel/account-service/internal/infra/helper"
)

func NewMessagingProducer(topicArn, awsEndpoint string) *MessagingProducer {
	return &MessagingProducer{
		TopicArn:    topicArn,
		AwsEndpoint: awsEndpoint,
	}
}

func (ref *MessagingProducer) Publish(eventType string, data interface{}, source string) error {
	log.Println("Publishing event", eventType, "from", source)
	awsConfig, err := helper.BuildAwsConfig(ref.AwsEndpoint)
	if err != nil {
		return err
	}
	snsClient := sns.NewFromConfig(awsConfig)

	dataAsJson, err := json.Marshal(data)
	if err != nil {
		message := "Error when convert event payload to json."
		log.Println(message, "Detail", err.Error())
		return exception.New(event.ErrorConvertMessageToJson, message, err, exception.HttpInternalError)
	}

	id := helper.NewULID()
	dataId := extractDataId(data)

	message := event.Event{
		Id:         id,
		OccurredAt: time.Now().UTC(),
		Type:       eventType,
		Source:     source,
		Data:       string(dataAsJson),
		DataId:     dataId,
	}
	messageAsJson, _ := json.Marshal(message)
	attrs := map[string]types.MessageAttributeValue{
		"EventType": {
			DataType:    aws.String("String"),
			StringValue: aws.String(eventType),
		},
		"Source": {
			DataType:    aws.String("String"),
			StringValue: aws.String(source),
		},
	}

	publishInput := &sns.PublishInput{
		Message:           aws.String(string(messageAsJson)),
		TopicArn:          aws.String(ref.TopicArn),
		MessageAttributes: attrs,
	}
	publishOutput, err := snsClient.Publish(context.TODO(), publishInput)
	if err != nil {
		message := "Error to publish event on topic."
		log.Println(message, "Detail", err.Error())
		return exception.New(event.ErrorConvertMessageToJson, message, err, exception.HttpInternalError)
	}

	log.Println(
		"Publishing event. Id", id, eventType, "from", source,
		"MessageId", *publishOutput.MessageId,
	)

	return nil
}

func (ref *MessagingConsumer) PollingMessages(messageChannel chan<- *sqsType.Message) {
	for {
		result, err := ref.SqsClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(ref.QueueUrl),
			MaxNumberOfMessages: ref.MaxNumberOfMessages,
			WaitTimeSeconds:     ref.WaitTimeSeconds,
			VisibilityTimeout:   ref.VisibilityTimeout,
		})

		if err != nil {
			log.Printf("Couldn't get messages from queue %v. Here's why: %v\n", ref.QueueUrl, err)
			continue
		}

		for _, message := range result.Messages {
			messageChannel <- &message
		}
	}
}

func (ref *MessagingConsumer) DeleteMessage(receiptHandle string) error {
	_, err := ref.SqsClient.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(ref.QueueUrl),
		ReceiptHandle: aws.String(receiptHandle),
	})

	if err != nil {
		log.Printf("Couldn't delete message from queue %v. Here's why: %v\n", ref.QueueUrl, err)
	}

	return err
}

func extractDataId(data interface{}) string {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	dataId := ""
	if v.Kind() == reflect.Struct {
		_, ok := v.Type().FieldByName("Id")
		if ok {
			return v.FieldByName("Id").Interface().(string)
		}
	}
	return dataId
}
