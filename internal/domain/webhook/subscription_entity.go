package webhook

import "time"

type Subscription struct {
	Id         string    `json:"id"`
	EventType  string    `json:"eventType"`
	Url        string    `json:"url"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	ExternalId string    `json:"externalId,omitempty"`
	Active     bool      `json:"active"`
}

// Error codes
const (
	SubscriptionNotFound = "SUBSCRIPTION_NOT_FOUND"
)
