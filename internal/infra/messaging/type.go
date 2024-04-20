package messaging

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/silasstoffel/account-service/configs"
	loggerContract "github.com/silasstoffel/account-service/internal/logger/contract"
)

type MessageSchema struct {
	Message           string
	MessageId         string
	MessageAttributes interface{}
}

type MessagingProducer struct {
	TopicArn    string
	AwsEndpoint string
	Config      *configs.Config
	Logger      loggerContract.Logger
}

type MessagingConsumer struct {
	SqsClient           *sqs.Client
	QueueUrl            string
	MaxNumberOfMessages int32
	WaitTimeSeconds     int32
	VisibilityTimeout   int32
	Logger              loggerContract.Logger
}

const (
	ErrorParserMessageFromTopic = "ErrorParserMessageFromTopic"
	ErrorParserMessageFromQueue = "ErrorParserMessageFromQueue"
)

func NewMessagingConsumer(queueUrl string, SqsClient *sqs.Client, logger loggerContract.Logger) *MessagingConsumer {
	return &MessagingConsumer{
		SqsClient:           SqsClient,
		QueueUrl:            queueUrl,
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     5,
		VisibilityTimeout:   30,
		Logger:              logger,
	}
}
