package database

import (
	"database/sql"

	"github.com/silasstoffel/account-service/internal/event"
	"github.com/silasstoffel/account-service/internal/exception"
	loggerContract "github.com/silasstoffel/account-service/internal/logger/contract"
)

type EventRepository struct {
	Db     *sql.DB
	Logger loggerContract.Logger
}

func NewEventRepository(db *sql.DB, logger loggerContract.Logger) *EventRepository {
	return &EventRepository{
		Db:     db,
		Logger: logger,
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
		repository.Logger.Error("Error when create event", err, nil)
		return exception.New(exception.DbCommandError, &err)
	}

	return nil
}
