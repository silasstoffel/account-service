package messaging

import "github.com/aws/aws-sdk-go-v2/service/sqs"

type MessageSchema struct {
	Message           string
	MessageId         string
	MessageAttributes interface{}
}

type MessagingProducer struct {
	TopicArn    string
	AwsEndpoint string
}

type MessagingConsumer struct {
	SqsClient           *sqs.Client
	QueueUrl            string
	MaxNumberOfMessages int32
	WaitTimeSeconds     int32
	VisibilityTimeout   int32
}

const (
	ErrorParserMessageFromTopic = "ErrorParserMessageFromTopic"
	ErrorParserMessageFromQueue = "ErrorParserMessageFromQueue"
)

func NewMessagingConsumer(queueUrl string, SqsClient *sqs.Client) *MessagingConsumer {
	return &MessagingConsumer{
		SqsClient:           SqsClient,
		QueueUrl:            queueUrl,
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     5,
		VisibilityTimeout:   30,
	}
}
