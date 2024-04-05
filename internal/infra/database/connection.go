package database

import (
	"database/sql"
	"fmt"
	"log"

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
		message := "Failed to open connection to database"
		log.Println(message, "Details:", err)
		return nil, exception.New(exception.UnknownError, message, err, 500)
	}

	err = db.Ping()
	if err != nil {
		message := "Failed to ping database"
		return nil, exception.New(exception.UnknownError, message, err, 500)
	}

	return db, nil
}

func CloseConnection(cnx *sql.DB) error {
	if err := cnx.Close(); err != nil {
		return exception.New(exception.UnknownError, "Failed to ping database", err, 500)
	}
	return nil
}
