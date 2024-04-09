package webhook

type CreateSubscriptionInput struct {
	Id         string `json:"id"`
	EventType  string `json:"eventType"`
	Url        string `json:"url"`
	ExternalId string `json:"externalId,omitempty"`
}

type SubscriptionRepository interface {
	GetByEventType(eventType string) ([]Subscription, error)
	Create(subscription CreateSubscriptionInput) (*Subscription, error)
}
