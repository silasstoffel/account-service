package database

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/silasstoffel/account-service/internal/domain/webhook"
	"github.com/silasstoffel/account-service/internal/exception"
	"github.com/silasstoffel/account-service/internal/infra/helper"
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
	var subscriptions []webhook.Subscription

	like := eventType
	before, _, found := strings.Cut(eventType, ".")
	if found {
		like = before + ".*"
	}

	stmt := `SELECT id, event_type, url, created_at, updated_at, external_id, active FROM webhook_subscriptions WHERE event_type IN($1, $2)`
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

	return subscriptions, nil
}

func (repository *SubscriptionRepository) Create(subscription webhook.CreateSubscriptionInput) (*webhook.Subscription, error) {
	lp := "[subscription-repository][create]"
	stmt := `INSERT INTO webhook_subscriptions (id, event_type, url, external_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`

	subscription.Id = helper.NewULID()
	now := time.Now().UTC()
	_, err := repository.Db.Exec(stmt,
		subscription.Id,
		subscription.EventType,
		subscription.Url,
		subscription.ExternalId,
		now,
		now,
	)

	if err != nil {
		log.Println(lp, "Error when creating webhook subscription", err.Error())
		return nil, exception.New(exception.UnknownError, "Error when creating webhook subscription", err, exception.HttpInternalError)
	}

	return &webhook.Subscription{
		Id:         subscription.Id,
		EventType:  subscription.EventType,
		Url:        subscription.Url,
		CreatedAt:  now,
		UpdatedAt:  now,
		ExternalId: subscription.ExternalId,
		Active:     true,
	}, nil
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
			&subscription.ExternalId,
			&subscription.Active,
		)
	case *sql.Rows:
		return r.Scan(
			&subscription.Id,
			&subscription.EventType,
			&subscription.Url,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
			&subscription.ExternalId,
			&subscription.Active,
		)
	}

	return exception.New(
		exception.UnknownError,
		"An Unknown error happens",
		errors.New("row argument is not sql.Row or sql.Rows"),
		exception.HttpInternalError,
	)
}
