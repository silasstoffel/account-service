package event

import "time"

// events
const (
	AccountCreated  = "account.created"
	AccountUpdated  = "account.updated"
	AccountInactive = "account.inactive"
)

// errors
const (
	ErrorPublishingEvent      = "ERROR_PUBLISHING_EVENT"
	ErrorConvertMessageToJson = "ERROR_CONVERT_MESSAGE_TO_JSON"
)

type Event struct {
	Id         string    `json:"id"`
	OccurredAt time.Time `json:"occurredAt"`
	Type       string    `json:"type"`
	Source     string    `json:"source"`
	Data       string    `json:"data"`
}

type EventService interface {
	Publish(eventType string, data interface{}, source string) error
}
