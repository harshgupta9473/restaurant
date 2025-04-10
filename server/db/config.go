package db

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	connStr := os.Getenv("connStr")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	DB = db
	return nil
}

func GetDB() *sql.DB {
	return DB
}
