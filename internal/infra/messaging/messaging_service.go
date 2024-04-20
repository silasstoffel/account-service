package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqsType "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	appConfig "github.com/silasstoffel/account-service/configs"
	"github.com/silasstoffel/account-service/internal/event"
	"github.com/silasstoffel/account-service/internal/exception"
	"github.com/silasstoffel/account-service/internal/infra/helper"
	loggerContract "github.com/silasstoffel/account-service/internal/logger/contract"
)

func NewMessagingProducer(topicArn, awsEndpoint string, config *appConfig.Config, logger loggerContract.Logger) *MessagingProducer {
	return &MessagingProducer{
		TopicArn:    topicArn,
		AwsEndpoint: awsEndpoint,
		Config:      config,
		Logger:      logger,
	}
}

func NewDefaultMessagingProducerFromConfig(config *appConfig.Config, logger loggerContract.Logger) *MessagingProducer {
	return NewMessagingProducer(
		config.Aws.AccountServiceTopicArn,
		config.Aws.Endpoint,
		config,
		logger,
	)
}

func (ref *MessagingProducer) Publish(eventType string, data interface{}, source string) error {
	prefix := "[messaging-service]"
	awsConfig, err := helper.BuildAwsConfig(ref.Config)
	if err != nil {
		ref.Logger.Error("Error when build aws config.", err, nil)
		return err
	}
	snsClient := sns.NewFromConfig(awsConfig)

	dataAsJson, err := json.Marshal(data)
	if err != nil {
		message := "Error when convert event payload to json."
		ref.Logger.Error(message, err, nil)
		return exception.New(exception.ErrorConvertMessageToJson, &err)
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
	if source == "" {
		source = "account-service"
	}
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
	_, err = snsClient.Publish(context.TODO(), publishInput)
	if err != nil {
		message := "Error to publish event on topic."
		ref.Logger.Error(fmt.Sprintf("%s %s", prefix, message), err, nil)
		return exception.New(exception.ErrorConvertMessageToJson, &err)
	}

	return nil
}

func (ref *MessagingConsumer) PollingMessages(messageChannel chan<- *sqsType.Message) {
	prefix := "[messaging-service]"
	for {
		result, err := ref.SqsClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(ref.QueueUrl),
			MaxNumberOfMessages: ref.MaxNumberOfMessages,
			WaitTimeSeconds:     ref.WaitTimeSeconds,
			VisibilityTimeout:   ref.VisibilityTimeout,
		})

		if err != nil {
			ref.Logger.Error(prefix+"Couldn't get messages from queue.", err, map[string]interface{}{
				"queueUrl": ref.QueueUrl,
			})
			continue
		}

		for _, message := range result.Messages {
			messageChannel <- &message
		}
	}
}

func (ref *MessagingConsumer) DeleteMessage(receiptHandle string) error {
	prefix := "[messaging-service]"
	_, err := ref.SqsClient.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(ref.QueueUrl),
		ReceiptHandle: aws.String(receiptHandle),
	})

	if err != nil {
		ref.Logger.Error(prefix+" Couldn't delete message from queue.", err, map[string]interface{}{
			"queueUrl": ref.QueueUrl,
		})
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
