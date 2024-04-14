package messaging

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/silasstoffel/account-service/internal/exception"
)

var topicMessageSchema MessageSchema

func ExtractMessageFromTopic(rawMessage *types.Message, message interface{}) error {
	if err := json.Unmarshal([]byte(*rawMessage.Body), &topicMessageSchema); err != nil {
		detail := "Error when parser message from topic"
		log.Println(message, detail, err)
		return exception.NewUnknown(&err)
	}

	if err := json.Unmarshal([]byte(topicMessageSchema.Message), &message); err != nil {
		detail := "Error when parser message from queue"
		log.Println(message, detail, err)
		return exception.NewUnknown(&err)
	}
	return nil
}

func ExtractMessageFromQueue(rawMessage *types.Message, message interface{}) error {
	if err := json.Unmarshal([]byte(*rawMessage.Body), &message); err != nil {
		detail := "Error when parser message from queue"
		log.Println(message, detail, err)
		return exception.NewUnknown(&err)
	}
	return nil
}
