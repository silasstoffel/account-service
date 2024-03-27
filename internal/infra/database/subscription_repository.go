package database

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/silasstoffel/account-service/internal/domain/webhook"
	"github.com/silasstoffel/account-service/internal/exception"
)

type SubscriptionRepository struct {
	Db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		Db: db,
	}
}

func (repository *SubscriptionRepository) GetByEventType(eventType string) ([]webhook.Subscription, error) {
	log.Println(loggerPrefix, "Finding subscriptions by event type")
	var subscriptions []webhook.Subscription

	like := eventType
	before, _, found := strings.Cut(eventType, ".")
	if found {
		like = before + ".*"
	}

	stmt := `SELECT id, event_type, url, created_at, updated_at FROM webhook_subscriptions WHERE event_type IN($1, $2)`
	rows, err := repository.Db.Query(stmt, eventType, like)
	if err != nil {
		log.Println(loggerPrefix, "error when execute command on database.", err.Error())
		return subscriptions, exception.New(exception.DbCommandError, "Error when listing subscriptions", err, exception.HttpInternalError)
	}

	message := "Error when scanning subscription"
	defer rows.Close()

	for rows.Next() {
		var subscription webhook.Subscription
		if err := scanSubscription(rows, &subscription); err != nil {
			log.Println(loggerPrefix, message, err.Error())
			return subscriptions, exception.New(exception.DbCommandError, message, err, exception.HttpInternalError)
		}
		subscriptions = append(subscriptions, subscription)
	}

	log.Println(loggerPrefix, "Finding subscription by event_type", eventType)

	return subscriptions, nil
}

func scanSubscription(row interface{}, subscription *webhook.Subscription) error {
	switch r := row.(type) {
	case *sql.Row:
		return r.Scan(
			&subscription.Id,
			&subscription.EventType,
			&subscription.Url,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		)
	case *sql.Rows:
		return r.Scan(
			&subscription.Id,
			&subscription.EventType,
			&subscription.Url,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		)
	}

	return exception.New(
		exception.UnknownError,
		"An Unknown error happens",
		errors.New("row argument is not sql.Row or sql.Rows"),
		exception.HttpInternalError,
	)
}
