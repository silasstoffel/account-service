package usecase

import (
	"github.com/silasstoffel/account-service/internal/domain/webhook"
	"github.com/silasstoffel/account-service/internal/event"
)

func (ref *WebHookSubscriptionUseCaseParams) ChangSubscriptionStatusUseCase(id string, active bool) error {
	subscription, err := ref.WebhookSubscriptionRepository.FindById(id)
	if err != nil {
		ref.Logger.Error("[change-subscription-status-usecase] Error when finding subscription", err, nil)
		return err
	}
	eventType := map[bool]string{
		true:  event.WebHookSubscriptionActivated,
		false: event.WebHookSubscriptionDeactivated,
	}
	toUpdate := webhook.UpdateSubscriptionInput{
		EventType:  subscription.EventType,
		Url:        subscription.Url,
		ExternalId: subscription.ExternalId,
		Active:     active,
	}
	subs, err := ref.WebhookSubscriptionRepository.Update(id, toUpdate)
	if err != nil {
		ref.Logger.Error("[change-subscription-status-usecase] Error when updating subscription", err, nil)
		return err
	}

	go ref.Messaging.Publish(eventType[active], subs, "account-service")
	go ref.Messaging.Publish(event.WebHookSubscriptionUpdated, subs, "account-service")

	return nil
}
