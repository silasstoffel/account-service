package webhook

import "time"

type UpdateTransactionInput struct {
	ReceivedStatusCode int       `json:"receivedStatusCode"`
	RequestStartedAt   time.Time `json:"requestStartedAt"`
	RequestFinishedAt  time.Time `json:"requestFinishedAt"`
}

type WebhookTransactionRepository interface {
	FindById(accountId string) (WebhookTransaction, error)
	Create(transaction WebhookTransaction) (WebhookTransaction, error)
	Update(id string, data UpdateTransactionInput) (WebhookTransaction, error)
}
