package database

import (
	"database/sql"
	"errors"
	"log"

	domain "github.com/silasstoffel/account-service/internal/domain/account"
	"github.com/silasstoffel/account-service/internal/exception"
)

type AccountPermissionRepository struct {
	Db *sql.DB
}

func NewAccountPermissionRepository(db *sql.DB) *AccountPermissionRepository {
	return &AccountPermissionRepository{
		Db: db,
	}
}

func (repository *AccountPermissionRepository) Create(data domain.AccountPermission) error {
	stmt := `INSERT INTO account_permissions(app_id, scope, account_id) VALUES ($1, $2, $3)`
	_, err := repository.Db.Exec(
		stmt,
		data.AppId,
		data.Scope,
		data.AccountId,
	)

	if err != nil {
		message := "Error when creating account permission"
		log.Println(message, "Detail:", err)
		return exception.New(exception.DbCommandError, message, err, exception.HttpInternalError)
	}

	return nil
}

func (repository *AccountPermissionRepository) DeleteByAccount(accountId string) error {
	stmt := `DELETE FROM account_permissions WHERE account_id = $1`
	_, err := repository.Db.Exec(stmt, accountId)

	if err != nil {
		message := "Error when deleting account permission"
		log.Println(message, "Detail:", err)
		return exception.New(exception.DbCommandError, message, err, exception.HttpInternalError)
	}

	return nil
}

func (repository *AccountPermissionRepository) FindByAccountId(accountId string) ([]domain.AccountPermission, error) {
	stmt := `SELECT account_id, app_id, scope FROM account_permissions WHERE account_id = $1`
	rows, err := repository.Db.Query(stmt, accountId)
	if err != nil {
		message := "Error when querying account permission"
		log.Println(message, "Detail:", err)
		return []domain.AccountPermission{}, exception.New(exception.DbCommandError, message, err, exception.HttpInternalError)
	}
	defer rows.Close()

	var result []domain.AccountPermission
	for rows.Next() {
		var data domain.AccountPermission
		err := scanAccountPermissionRow(rows, &data)
		if err != nil {
			return []domain.AccountPermission{}, err
		}
		result = append(result, data)
	}

	return result, nil
}

func scanAccountPermissionRow(row interface{}, data *domain.AccountPermission) error {
	switch r := row.(type) {
	case *sql.Row:
		return r.Scan(
			&data.AccountId,
			&data.AppId,
			&data.Scope,
		)
	case *sql.Rows:
		return r.Scan(
			&data.AccountId,
			&data.AppId,
			&data.Scope,
		)
	}
	return exception.New(
		exception.UnknownError,
		"An Unknown error happens",
		errors.New("ScanRow error is not sql.Row or sql.Rows"),
		exception.HttpInternalError,
	)
}
