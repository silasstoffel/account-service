package database

import (
	"database/sql"
	"errors"
	"fmt"

	domain "github.com/silasstoffel/account-service/internal/domain/account"
	"github.com/silasstoffel/account-service/internal/exception"
	loggerContract "github.com/silasstoffel/account-service/internal/logger/contract"
)

type PermissionRepository struct {
	Db     *sql.DB
	Logger loggerContract.Logger
}

func NewPermissionRepository(db *sql.DB, logger loggerContract.Logger) *PermissionRepository {
	return &PermissionRepository{
		Db:     db,
		Logger: logger,
	}
}

func buildReadPermissionCommand(where, orderBy, limit, offset string) string {
	if where == "" {
		where = "1=1"
	}
	if orderBy == "" {
		orderBy = "1"
	}

	stmt := `SELECT
				id,
				scope,
				active,
				created_at
			FROM permissions
			WHERE %s
			ORDER BY %s`

	query := fmt.Sprintf(stmt, where, orderBy)

	if limit != "" && offset != "" {
		query = fmt.Sprintf("%s LIMIT %s OFFSET %s", query, limit, offset)
	}

	return query
}

func (repository *PermissionRepository) List(input domain.ListPermissionInput) ([]*domain.Permission, error) {
	offset, limit, _ := buildPaginationParams(input.Limit, input.Page)
	query := buildReadPermissionCommand("", "", "$1", "$2")
	rows, err := repository.Db.Query(query, limit, offset)
	lp := "[permission-repository][list]"

	if err != nil {
		message := lp + " Failure when querying permission"
		repository.Logger.Error(message, err, nil)
		return nil, exception.NewUnknownError(&err)
	}
	defer rows.Close()

	var result []*domain.Permission
	for rows.Next() {
		var data domain.Permission
		err := scanPermissionRow(rows, &data)
		if err != nil {
			repository.Logger.Error(lp+" Failure when scan permission", err, nil)
			return nil, err
		}
		result = append(result, &data)
	}

	return result, nil
}

func scanPermissionRow(row interface{}, data *domain.Permission) error {
	switch r := row.(type) {
	case *sql.Row:
		return r.Scan(
			&data.Id,
			&data.Scope,
			&data.Active,
			&data.CreatedAt,
		)
	case *sql.Rows:
		return r.Scan(
			&data.Id,
			&data.Scope,
			&data.Active,
			&data.CreatedAt,
		)
	}
	e := errors.New("ScanRow error is not sql.Row or sql.Rows")
	return exception.NewUnknownError(&e)
}
