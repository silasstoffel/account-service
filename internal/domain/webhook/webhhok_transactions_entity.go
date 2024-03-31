package webhook

import "time"

type WebhookTransaction struct {
	Id                 string    `json:"id"`
	EventId            string    `json:"eventId"`
	SubscriptionId     string    `json:"subscriptionId"`
	EventType          string    `json:"eventType"`
	ReceivedStatusCode int       `json:"receivedStatusCode"`
	RequestStartedAt   time.Time `json:"requestStartedAt"`
	RequestFinishedAt  time.Time `json:"requestFinishedAt"`
	NumberOfRequests   int       `json:"numberOfRequests"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

// Error codes
const (
	WebhookTransactionNotFound            = "WEBHOOK_TRANSACTION_NOT_FOUND"
	WebhookTransactionNotificationTimeout = "WEBHOOK_TRANSACTION_NOTIFICATION_TIMEOUT"
)
