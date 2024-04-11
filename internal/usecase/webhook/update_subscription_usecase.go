package usecase

import (
	"log"

	"github.com/silasstoffel/account-service/internal/domain/webhook"
	"github.com/silasstoffel/account-service/internal/event"
)

func (ref *WebHookSubscriptionUseCaseParams) UpdateSubscriptionUseCase(id string, input webhook.UpdateSubscriptionInput) (*webhook.Subscription, error) {
	subscription, err := ref.WebhookSubscriptionRepository.FindById(id)
	if err != nil {
		log.Println("Error when updating subscription", err)
		return nil, err
	}

	subs, err := ref.WebhookSubscriptionRepository.Update(id, buildInput(subscription, input))
	if err != nil {
		log.Println("Error when updating subscription", err)
		return nil, err
	}
	subs.CreatedAt = subscription.CreatedAt

	go ref.Messaging.Publish(event.WebHookSubscriptionUpdated, subs, "account-service")

	return subs, nil
}

func buildInput(subscription *webhook.Subscription, input webhook.UpdateSubscriptionInput) webhook.UpdateSubscriptionInput {
	if input.EventType == "" {
		input.EventType = subscription.EventType
	}
	if input.Url == "" {
		input.Url = subscription.Url
	}
	if input.ExternalId == "" {
		input.ExternalId = subscription.ExternalId
	}
	return input
}
