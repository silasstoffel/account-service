package usecase

import (
	"github.com/silasstoffel/account-service/internal/domain/webhook"
	"github.com/silasstoffel/account-service/internal/event"
	loggerContract "github.com/silasstoffel/account-service/internal/logger/contract"
)

type WebHookSubscriptionUseCaseParams struct {
	Messaging                     event.EventProducer
	WebhookSubscriptionRepository webhook.SubscriptionRepository
	Logger                        loggerContract.Logger
}
