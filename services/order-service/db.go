package main

import (
	"database/sql"
	_ "github.com/lib/pq" // or use mysql driver
	"log"
)

func InitDB() *sql.DB {
	connStr := "user=postgres password=yourpassword dbname=yourdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}
