package usecase

import (
	"github.com/silasstoffel/account-service/internal/domain/webhook"
	"github.com/silasstoffel/account-service/internal/event"
)

type WebHookSubscriptionUseCaseParams struct {
	Messaging                     event.EventProducer
	WebhookSubscriptionRepository webhook.SubscriptionRepository
}
