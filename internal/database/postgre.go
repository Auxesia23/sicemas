package database

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitDB() (*sqlx.DB, error) {
	uri := os.Getenv("DB_URI")

	db, err := sqlx.Open("postgres", uri)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)

	log.Println("Postgre Connected Succesfuly")
	return db, nil
}
