package webhook

type SubscriptionRepository interface {
	GetByEventType(eventType string) ([]Subscription, error)
}
