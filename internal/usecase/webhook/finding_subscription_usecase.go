package usecase

import (
	"github.com/silasstoffel/account-service/internal/domain/webhook"
)

func (ref *WebHookSubscriptionUseCaseParams) FindSubscriptionUseCase(id string) (*webhook.Subscription, error) {
	subscription, err := ref.WebhookSubscriptionRepository.FindById(id)
	if err != nil {
		ref.Logger.Error("[find-subscription-usecase] Error when finding subscription", err, nil)
		return nil, err
	}
	return subscription, nil
}
