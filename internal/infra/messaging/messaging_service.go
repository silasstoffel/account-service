package messaging

import (
	"encoding/json"
	"log"
	"time"

	"github.com/silasstoffel/account-service/internal/event"
	"github.com/silasstoffel/account-service/internal/exception"
	"github.com/silasstoffel/account-service/internal/infra/helper"
)

type MessagingService struct {
	Type string
}

func NewMessagingService() *MessagingService {
	return &MessagingService{
		Type: "sqs",
	}
}

func (ref *MessagingService) Publish(eventType string, data interface{}, source string) error {
	log.Println("Publishing event", eventType, "from", source)

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

	log.Println("Publishing event", "id", id, eventType, "from", source)
	log.Println("Event data", message)

	return nil
}
