package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	domain "github.com/silasstoffel/account-service/internal/domain/account"
	"github.com/silasstoffel/account-service/internal/exception"
)

type PermissionRepository struct {
	Db *sql.DB
}

func NewPermissionRepository(db *sql.DB) *PermissionRepository {
	return &PermissionRepository{
		Db: db,
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
		message := "Failure when querying permission"
		log.Println(lp, message, err)
		return nil, exception.NewUnknown(&err)
	}
	defer rows.Close()

	var result []*domain.Permission
	for rows.Next() {
		var data domain.Permission
		err := scanPermissionRow(rows, &data)
		if err != nil {
			log.Println(lp, "Failure when scan permission", err)
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
	return exception.NewUnknown(&e)
}
