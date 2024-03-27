package usecase

import (
	"log"

	"github.com/silasstoffel/account-service/internal/event"
)

type CreateEventUseCase struct {
	EventRepository event.EventRepository
	Messaging       event.EventProducer
}

func (ref *CreateEventUseCase) CreateEventUseCase(input event.Event) error {
	loggerPrefix := "[create-event-usecase]"
	log.Println(loggerPrefix, "Creating event - Type:", input.Type)

	err := ref.EventRepository.Create(input)

	if err != nil {
		return err
	}

	log.Println(loggerPrefix, "Event created", "id:", input.Id)

	go ref.Messaging.Publish(event.EventCreated, input, "account-service")

	return nil
}
