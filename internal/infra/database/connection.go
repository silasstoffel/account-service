package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const prefix = "[database]"

func OpenConnection() *sql.DB {
	log.Println(prefix, "Opening postgres connection...")

	db, err := sql.Open("postgres", "user=account password=account dbname=account-service sslmode=disable")

	if err != nil {
		log.Println(prefix, "error connection")
		panic("Error on create connection.")
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
