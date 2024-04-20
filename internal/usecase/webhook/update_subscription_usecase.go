package usecase

import (
	"github.com/silasstoffel/account-service/internal/domain/webhook"
	"github.com/silasstoffel/account-service/internal/event"
)

func (ref *WebHookSubscriptionUseCaseParams) UpdateSubscriptionUseCase(id string, input webhook.UpdateSubscriptionInput) (*webhook.Subscription, error) {
	subscription, err := ref.WebhookSubscriptionRepository.FindById(id)
	if err != nil {
		ref.Logger.Error("[update-subscription-usecase] Error when finding subscription", err, nil)
		return nil, err
	}

	subs, err := ref.WebhookSubscriptionRepository.Update(id, buildInput(subscription, input))
	if err != nil {
		ref.Logger.Error("[update-subscription-usecase] Error when updating subscription", err, nil)
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
