package webhook

type CreateSubscriptionInput struct {
	Id         string `json:"id"`
	EventType  string `json:"eventType" validate:"required"`
	Url        string `json:"url" validate:"required"`
	ExternalId string `json:"externalId,omitempty" validate:"required"`
}

type UpdateSubscriptionInput struct {
	EventType  string `json:"eventType"`
	Url        string `json:"url"`
	ExternalId string `json:"externalId,omitempty"`
	Active     bool   `json:"active,omitempty"`
}

type ListSubscriptionInput struct {
	Page  int
	Limit int
}

type SubscriptionReadRepository interface {
	GetByEventType(eventType string) ([]Subscription, error)
	FindById(id string) (*Subscription, error)
	List(input ListSubscriptionInput) ([]*Subscription, error)
}

type SubscriptionWriteRepository interface {
	Create(subscription CreateSubscriptionInput) (*Subscription, error)
	Update(id string, data UpdateSubscriptionInput) (*Subscription, error)
}

type SubscriptionRepository interface {
	SubscriptionReadRepository
	SubscriptionWriteRepository
}
