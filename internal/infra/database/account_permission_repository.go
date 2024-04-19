package database

import (
	"database/sql"
	"errors"
	"fmt"

	domain "github.com/silasstoffel/account-service/internal/domain/account"
	"github.com/silasstoffel/account-service/internal/exception"
	loggerContract "github.com/silasstoffel/account-service/internal/logger/contract"
)

type AccountPermissionRepository struct {
	Db     *sql.DB
	Logger loggerContract.Logger
}

func NewAccountPermissionRepository(db *sql.DB, logger loggerContract.Logger) *AccountPermissionRepository {
	return &AccountPermissionRepository{
		Db:     db,
		Logger: logger,
	}
}

func buildAccountPermissionSelectCommand(where, orderBy, limit, offset string) string {
	if where == "" {
		where = "1=1"
	}
	if orderBy == "" {
		orderBy = "1"
	}

	stmt := `SELECT
				ap.account_id,
				ap.permission_id,
				ap.created_at,
				p.scope,
				p.active
			FROM account_permissions ap
				 JOIN permissions p ON p.id = ap.permission_id
			WHERE %s
			ORDER BY %s`

	query := fmt.Sprintf(stmt, where, orderBy)

	if limit != "" && offset != "" {
		query = fmt.Sprintf("%s LIMIT %s OFFSET %s", query, limit, offset)
	}

	return query
}

func (repository *AccountPermissionRepository) Create(data domain.CreateAccountPermissionInput) error {
	stmt := `INSERT INTO account_permissions(account_id, permission_id) VALUES ($1, $2)`
	_, err := repository.Db.Exec(
		stmt,
		data.AccountId,
		data.PermissionId,
	)

	if err != nil {
		message := "Error when creating account permission"
		repository.Logger.Error(message, err, nil)
		return exception.New(exception.DbCommandError, &err)
	}

	return nil
}

func (repository *AccountPermissionRepository) DeleteByAccount(accountId string) error {
	stmt := `DELETE FROM account_permissions WHERE account_id = $1`
	_, err := repository.Db.Exec(stmt, accountId)

	if err != nil {
		message := "Error when deleting account permission"
		repository.Logger.Error(message, err, nil)
		return exception.New(exception.DbCommandError, &err)
	}

	return nil
}

func (repository *AccountPermissionRepository) FindByAccountId(accountId string) ([]domain.AccountPermission, error) {
	query := buildAccountPermissionSelectCommand("ap.account_id = $1", "", "", "")
	rows, err := repository.Db.Query(query, accountId)
	if err != nil {
		message := "Error when querying account permission"
		repository.Logger.Error(message, err, nil)
		return nil, exception.New(exception.DbCommandError, &err)
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
			&data.PermissionId,
			&data.CreatedAt,
			&data.Scope,
			&data.Active,
		)
	case *sql.Rows:
		return r.Scan(
			&data.AccountId,
			&data.PermissionId,
			&data.CreatedAt,
			&data.Scope,
			&data.Active,
		)
	}
	e := errors.New("ScanRow error is not sql.Row or sql.Rows")
	return exception.New(exception.UnknownError, &e)
}
