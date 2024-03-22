package database

import (
	"database/sql"
	"log"

	"github.com/silasstoffel/account-service/internal/event"
	"github.com/silasstoffel/account-service/internal/exception"
)

type EventRepository struct {
	Db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{
		Db: db,
	}
}

func (repository *EventRepository) Create(event event.Event) error {
	loggerPrefix := "[event-repository]"
	log.Println(loggerPrefix, "Creating event...")

	stmt := `INSERT INTO events (id, occurred_at, "type", "source", "data") VALUES ($1, $2, $3, $4, $5)`
	_, err := repository.Db.Exec(
		stmt,
		event.Id,
		event.OccurredAt,
		event.Type,
		event.Source,
		event.Data,
	)

	if err != nil {
		log.Println(loggerPrefix, "Error when create event", err.Error())
		return exception.New(exception.DbCommandError, "Error when creating event", err, exception.HttpInternalError)
	}

	log.Println(loggerPrefix, "Event created with id", event.Id)
	return nil
}
