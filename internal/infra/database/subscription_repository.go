package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/silasstoffel/account-service/internal/domain/webhook"
	"github.com/silasstoffel/account-service/internal/exception"
	"github.com/silasstoffel/account-service/internal/infra/helper"
	loggerContract "github.com/silasstoffel/account-service/internal/logger/contract"
)

type SubscriptionRepository struct {
	Db     *sql.DB
	Logger loggerContract.Logger
}

func NewSubscriptionRepository(db *sql.DB, Logger loggerContract.Logger) *SubscriptionRepository {
	return &SubscriptionRepository{
		Db:     db,
		Logger: Logger,
	}
}

func buildSelectCommand(where, orderBy, limit, offset string) string {
	if where == "" {
		where = "1=1"
	}
	if orderBy == "" {
		orderBy = "1"
	}
	query := fmt.Sprintf("SELECT id, event_type, url, created_at, updated_at, external_id, active FROM webhook_subscriptions WHERE %s ORDER BY %s", where, orderBy)

	if limit != "" && offset != "" {
		query = fmt.Sprintf("%s LIMIT %s OFFSET %s", query, limit, offset)
	}

	return query
}

func (repository *SubscriptionRepository) GetByEventType(eventType string) ([]webhook.Subscription, error) {
	var subscriptions []webhook.Subscription

	like := eventType
	before, _, found := strings.Cut(eventType, ".")
	if found {
		like = before + ".*"
	}

	stmt := buildSelectCommand("event_type IN($1, $2)", "id", "", "")
	rows, err := repository.Db.Query(stmt, eventType, like)
	if err != nil {
		repository.Logger.Error("Error when execute command on database.", err, nil)
		return subscriptions, exception.NewDbCommandError(&err)
	}

	message := "Error when scanning subscription"
	defer rows.Close()

	for rows.Next() {
		var subscription webhook.Subscription
		if err := scanSubscription(rows, &subscription); err != nil {
			repository.Logger.Error(message, err, nil)
			return subscriptions, exception.NewDbCommandError(&err)
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
		repository.Logger.Error(lp+" Error when creating webhook subscription", err, nil)
		return nil, exception.NewUnknownError(&err)
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

func (repository *SubscriptionRepository) FindById(id string) (*webhook.Subscription, error) {
	stmt := buildSelectCommand("id = $1", "id", "", "")

	var subscription webhook.Subscription
	row := repository.Db.QueryRow(stmt, id)

	if err := scanSubscription(row, &subscription); err != nil {
		lp := "[subscription-repository][get-by-id]"
		message := lp + " Error when finding subscription"
		if err == sql.ErrNoRows {
			return nil, exception.New(exception.WebhookSubscriptionNotFound, &err)
		}
		repository.Logger.Error(message, err, nil)
		return nil, exception.NewUnknownError(&err)
	}

	return &subscription, nil
}

func (repository *SubscriptionRepository) Update(id string, data webhook.UpdateSubscriptionInput) (*webhook.Subscription, error) {
	stmt := `UPDATE webhook_subscriptions SET event_type = $1, url = $2, external_id = $3, updated_at = $4, active = $5 WHERE id = $6`

	now := time.Now().UTC()
	_, err := repository.Db.Exec(stmt,
		data.EventType,
		data.Url,
		data.ExternalId,
		now,
		data.Active,
		id,
	)

	if err != nil {
		repository.Logger.Error("[subscription-repository][update] Error when updating webhook subscription", err, nil)
		return nil, exception.NewUnknownError(&err)
	}

	return &webhook.Subscription{
		Id:         id,
		EventType:  data.EventType,
		Url:        data.Url,
		UpdatedAt:  now,
		ExternalId: data.ExternalId,
		Active:     data.Active,
	}, nil
}

func (repository *SubscriptionRepository) List(input webhook.ListSubscriptionInput) ([]*webhook.Subscription, error) {
	page := input.Page
	limit := input.Limit

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 12
	}

	offset := (page - 1) * limit
	stmt := buildSelectCommand("", "id", "$1", "$2")
	rows, err := repository.Db.Query(stmt, limit, offset)
	lp := "[subscription-repository][list]"
	if err != nil {
		message := lp + " Error when find subscriptions"
		repository.Logger.Error(message, err, nil)
		return nil, exception.NewDbCommandError(&err)
	}

	var subscriptions []*webhook.Subscription
	defer rows.Close()

	for rows.Next() {
		var sub webhook.Subscription
		if err := scanSubscription(rows, &sub); err != nil {
			repository.Logger.Error(lp+" Error when scan result", err, nil)
			return nil, exception.NewDbCommandError(&err)
		}
		subscriptions = append(subscriptions, &sub)
	}

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
	e := errors.New("row argument is not sql.Row or sql.Rows")
	return exception.NewUnknownError(&e)
}
