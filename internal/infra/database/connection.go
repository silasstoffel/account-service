package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/silasstoffel/account-service/configs"
	"github.com/silasstoffel/account-service/internal/exception"
)

func OpenConnection(config *configs.Config) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Db.Host,
		config.Db.Port,
		config.Db.User,
		config.Db.Password,
		config.Db.Name,
	)

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, exception.New(exception.UnknownError, &err)
	}

	err = db.Ping()
	if err != nil {
		return nil, exception.New(exception.UnknownError, &err)
	}

	return db, nil
}

func CloseConnection(cnx *sql.DB) error {
	if err := cnx.Close(); err != nil {
		return exception.New(exception.UnknownError, &err)
	}
	return nil
}
