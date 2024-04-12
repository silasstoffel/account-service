package usecase

import (
	"log"

	"github.com/silasstoffel/account-service/internal/domain/webhook"
)

func (ref *WebHookSubscriptionUseCaseParams) ListSubscriptionUseCase(input webhook.ListSubscriptionInput) ([]*webhook.Subscription, error) {
	subscriptions, err := ref.WebhookSubscriptionRepository.List(input)
	if err != nil {
		log.Println("Error when list subscription", err)
		return nil, err
	}
	return subscriptions, nil
}
