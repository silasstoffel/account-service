package event

import "time"

// events
const (
	AccountCreated                 = "account.created"
	AccountUpdated                 = "account.updated"
	AccountInactive                = "account.inactive"
	AccountLogged                  = "account.logged"
	EventCreated                   = "event.created"
	WebHookSubscriptionCreated     = "webhook.subscription.created"
	WebHookSubscriptionUpdated     = "webhook.subscription.updated"
	WebHookSubscriptionActivated   = "webhook.subscription.activated"
	WebHookSubscriptionDeactivated = "webhook.subscription.deactivated"
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
