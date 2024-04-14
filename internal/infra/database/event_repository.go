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
	stmt := `INSERT INTO events (id, occurred_at, "type", "source", "data", data_id) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := repository.Db.Exec(
		stmt,
		event.Id,
		event.OccurredAt,
		event.Type,
		event.Source,
		event.Data,
		event.DataId,
	)

	if err != nil {
		lp := "[event-repository]"
		log.Println(lp, "Error when create event", err.Error())
		return exception.New(exception.DbCommandError, &err)
	}

	return nil
}
