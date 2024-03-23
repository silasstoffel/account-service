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
}

const (
	ErrorParserMessageFromTopic = "ErrorParserMessageFromTopic"
	ErrorParserMessageFromQueue = "ErrorParserMessageFromQueue"
)
