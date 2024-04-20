package usecase

import (
	"github.com/silasstoffel/account-service/internal/domain/webhook"
	"github.com/silasstoffel/account-service/internal/event"
)

func (ref *WebHookSubscriptionUseCaseParams) CreateSubscriptionUseCase(input webhook.CreateSubscriptionInput) (*webhook.Subscription, error) {
	subs, err := ref.WebhookSubscriptionRepository.Create(input)
	if err != nil {
		ref.Logger.Error("[create-subscription-usecase] Error when creating subscription", err, nil)
		return nil, err
	}

	go ref.Messaging.Publish(event.WebHookSubscriptionCreated, subs, "account-service")

	return subs, nil
}
