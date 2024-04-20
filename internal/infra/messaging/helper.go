package messaging

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/silasstoffel/account-service/internal/exception"
)

var topicMessageSchema MessageSchema

func ExtractMessageFromTopic(rawMessage *types.Message, message interface{}) error {
	if err := json.Unmarshal([]byte(*rawMessage.Body), &topicMessageSchema); err != nil {
		return exception.NewUnknownError(&err)
	}

	if err := json.Unmarshal([]byte(topicMessageSchema.Message), &message); err != nil {
		return exception.NewUnknownError(&err)
	}
	return nil
}

func ExtractMessageFromQueue(rawMessage *types.Message, message interface{}) error {
	if err := json.Unmarshal([]byte(*rawMessage.Body), &message); err != nil {
		return exception.NewUnknownError(&err)
	}
	return nil
}
