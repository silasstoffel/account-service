package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/silasstoffel/account-service/configs"
)

const prefix = "[database]"

func OpenConnection(config *configs.Config) *sql.DB {
	log.Println(prefix, "Opening postgres connection...")
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
		log.Println(prefix, "error connection")
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Println(prefix, "Error when ping database")
		panic(err)
	}
	log.Println(prefix, "Connection opened")
	return db
}

func CloseConnection(cnx *sql.DB) {
	log.Println(prefix, "Closing connection")
	cnx.Close()
	log.Println(prefix, "Connection closed")
}
