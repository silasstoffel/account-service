package database

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/silasstoffel/account-service/internal/domain/webhook"
	"github.com/silasstoffel/account-service/internal/exception"
)

type WebhookTransactionRepository struct {
	Db *sql.DB
}

func NewWebhookTransactionRepository(db *sql.DB) *WebhookTransactionRepository {
	return &WebhookTransactionRepository{
		Db: db,
	}
}

func (repository *WebhookTransactionRepository) FindById(id string) (webhook.WebhookTransaction, error) {
	var transaction webhook.WebhookTransaction

	stmt := `SELECT id, event_id, subscription_id, event_type,
		received_status_code, started_at, finished_at,
		number_of_requests, created_at, updated_at
	FROM webhook_transactions
	WHERE id = $1`

	row := repository.Db.QueryRow(stmt, id)
	if err := scanTransactions(row, &transaction); err != nil {
		log.Println("Error when finding transactions by id.", err)
		if err == sql.ErrNoRows {
			return transaction, exception.New(webhook.WebhookTransactionNotFound, &err)
		}
		return transaction, exception.NewDbCommandError(&err)
	}

	return transaction, nil
}

func (repository *WebhookTransactionRepository) Create(transaction webhook.WebhookTransaction) (webhook.WebhookTransaction, error) {
	transaction.CreatedAt = time.Now().UTC()
	transaction.UpdatedAt = transaction.CreatedAt
	stmt := `INSERT INTO webhook_transactions(
				id,
				event_id,
				subscription_id,
				event_type,
				received_status_code,
				started_at,
				finished_at,
				number_of_requests,
				created_at,
				updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := repository.Db.Exec(
		stmt,
		transaction.Id,
		transaction.EventId,
		transaction.SubscriptionId,
		transaction.EventType,
		transaction.ReceivedStatusCode,
		transaction.RequestStartedAt,
		transaction.RequestFinishedAt,
		transaction.NumberOfRequests,
		transaction.CreatedAt,
		transaction.UpdatedAt,
	)

	if err != nil {
		message := "Error when creating webhook transaction"
		log.Println(message, err)
		return transaction, exception.NewDbCommandError(&err)
	}

	return transaction, nil
}

func (repository *WebhookTransactionRepository) Update(id string, transaction webhook.UpdateTransactionInput) (webhook.WebhookTransaction, error) {
	toUpdate, err := (repository).FindById(id)
	if err != nil {
		log.Println("Error when finding transaction", id, "Detail", err)
		return toUpdate, err
	}

	stmt := `UPDATE webhook_transactions
			SET updated_at = $1,
				received_status_code = $2,
				started_at = $3,
				finished_at = $4,
				number_of_requests = number_of_requests + 1
			WHERE id = $5`

	now := time.Now().UTC()
	_, err = repository.Db.Exec(
		stmt,
		now,
		transaction.ReceivedStatusCode,
		transaction.RequestStartedAt,
		transaction.RequestFinishedAt,
		id,
	)

	if err != nil {
		message := "Error when updating webhook transactions"
		log.Println(loggerPrefix, message, "id:", id, "eventId:", toUpdate.EventId, err)
		return toUpdate, exception.NewDbCommandError(&err)
	}

	toUpdate.ReceivedStatusCode = transaction.ReceivedStatusCode
	toUpdate.RequestStartedAt = transaction.RequestStartedAt
	toUpdate.RequestFinishedAt = transaction.RequestFinishedAt
	toUpdate.UpdatedAt = now
	toUpdate.NumberOfRequests++

	return toUpdate, nil
}

func scanTransactions(row interface{}, transaction *webhook.WebhookTransaction) error {
	switch r := row.(type) {
	case *sql.Row:
		return r.Scan(
			&transaction.Id,
			&transaction.EventId,
			&transaction.SubscriptionId,
			&transaction.EventType,
			&transaction.ReceivedStatusCode,
			&transaction.RequestStartedAt,
			&transaction.RequestFinishedAt,
			&transaction.NumberOfRequests,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
	case *sql.Rows:
		return r.Scan(
			&transaction.Id,
			&transaction.EventId,
			&transaction.SubscriptionId,
			&transaction.EventType,
			&transaction.ReceivedStatusCode,
			&transaction.RequestStartedAt,
			&transaction.RequestFinishedAt,
			&transaction.NumberOfRequests,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
	}
	e := errors.New("row argument is not sql.Row or sql.Rows")
	return exception.NewUnknownError(&e)
}
