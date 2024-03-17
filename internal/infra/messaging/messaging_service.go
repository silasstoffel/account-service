package messaging

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/silasstoffel/account-service/internal/event"
	"github.com/silasstoffel/account-service/internal/exception"
	"github.com/silasstoffel/account-service/internal/infra/helper"
)

type MessagingService struct {
	TopicArn    string
	AwsEndpoint string
}

func NewMessagingService(topicArn, awsEndpoint string) *MessagingService {
	return &MessagingService{
		TopicArn:    topicArn,
		AwsEndpoint: awsEndpoint,
	}
}

func (ref *MessagingService) Publish(eventType string, data interface{}, source string) error {
	log.Println("Publishing event", eventType, "from", source)
	awsConfig, err := buildAwsConfig(ref.AwsEndpoint)
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
	message := event.Event{
		Id:         id,
		OccurredAt: time.Now().UTC(),
		Type:       eventType,
		Source:     source,
		Data:       string(dataAsJson),
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
		message := "Error when convert event payload to json."
		log.Println(message, "Detail", err.Error())
		return exception.New(event.ErrorConvertMessageToJson, message, err, exception.HttpInternalError)
	}

	log.Println(
		"Publishing event. Id", id, eventType, "from", source,
		"MessageId", *publishOutput.MessageId,
	)

	return nil
}

func buildAwsConfig(awsEndpoint string) (cfg aws.Config, err error) {
	awsRegion := "us-east-1"
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if awsEndpoint != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           awsEndpoint,
				SigningRegion: awsRegion,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithEndpointResolverWithOptions(customResolver),
	)

	if err != nil {
		return aws.Config{}, exception.New(event.ErrorInstanceEventBus, "Error creating event bus instance", err, exception.HttpInternalError)
	}

	return awsCfg, nil
}
