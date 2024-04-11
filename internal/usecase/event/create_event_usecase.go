package usecase

import (
	"github.com/silasstoffel/account-service/internal/event"
)

type CreateEventUseCase struct {
	EventRepository event.EventRepository
	Messaging       event.EventProducer
}

func (ref *CreateEventUseCase) CreateEventUseCase(input event.Event) error {
	err := ref.EventRepository.Create(input)
	if err != nil {
		return err
	}
	go ref.Messaging.Publish(event.EventCreated, input, "account-service")

	return nil
}
