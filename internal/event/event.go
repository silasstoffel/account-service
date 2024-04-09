package event

import "time"

// events
const (
	AccountCreated             = "account.created"
	AccountUpdated             = "account.updated"
	AccountInactive            = "account.inactive"
	AccountLogged              = "account.logged"
	EventCreated               = "event.created"
	WebHookSubscriptionCreated = "webhook.subscription.created"
	WebHookSubscriptionUpdated = "webhook.subscription.updated"
)

// errors
const (
	ErrorPublishingEvent      = "ERROR_PUBLISHING_EVENT"
	ErrorConvertMessageToJson = "ERROR_CONVERT_MESSAGE_TO_JSON"
	ErrorInstanceEventBus     = "ERROR_INSTANCE_EVENT_BUS"
)

type Event struct {
	Id         string    `json:"id"`
	DataId     string    `json:"dataId"`
	OccurredAt time.Time `json:"occurredAt"`
	Type       string    `json:"type"`
	Source     string    `json:"source"`
	Data       string    `json:"data"`
}

type EventProducer interface {
	Publish(eventType string, data interface{}, source string) error
}

type EventRepository interface {
	Create(event Event) error
}
