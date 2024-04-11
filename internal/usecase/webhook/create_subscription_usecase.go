package usecase

import (
	"log"

	"github.com/silasstoffel/account-service/internal/domain/webhook"
	"github.com/silasstoffel/account-service/internal/event"
)

func (ref *WebHookSubscriptionUseCaseParams) CreateSubscriptionUseCase(input webhook.CreateSubscriptionInput) (*webhook.Subscription, error) {
	subs, err := ref.WebhookSubscriptionRepository.Create(input)
	if err != nil {
		log.Println("Error when create subscription", err)
		return nil, err
	}

	go ref.Messaging.Publish(event.WebHookSubscriptionCreated, subs, "account-service")

	return subs, nil
}
