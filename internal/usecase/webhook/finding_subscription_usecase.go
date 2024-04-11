package usecase

import (
	"log"

	"github.com/silasstoffel/account-service/internal/domain/webhook"
)

func (ref *WebHookSubscriptionUseCaseParams) FindSubscriptionUseCase(id string) (*webhook.Subscription, error) {
	subscription, err := ref.WebhookSubscriptionRepository.FindById(id)
	if err != nil {
		log.Println("Error when updating subscription", err)
		return nil, err
	}
	return subscription, nil
}
