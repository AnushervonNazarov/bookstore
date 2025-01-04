package db

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func InitDB(connStr string) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	DB = db
}
