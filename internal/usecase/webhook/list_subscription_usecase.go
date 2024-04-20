package usecase

import (
	"github.com/silasstoffel/account-service/internal/domain/webhook"
)

func (ref *WebHookSubscriptionUseCaseParams) ListSubscriptionUseCase(input webhook.ListSubscriptionInput) ([]*webhook.Subscription, error) {
	subscriptions, err := ref.WebhookSubscriptionRepository.List(input)
	if err != nil {
		ref.Logger.Error("[list-subscription-usecase] Error when listing subscriptions", err, nil)
		return nil, err
	}
	return subscriptions, nil
}
