package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq" // or use mysql driver
)

func InitDB() *sql.DB {
	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("failed to open DB:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("failed to ping DB:", err)
	}
	return db
}
