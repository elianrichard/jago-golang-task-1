package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	// Open database with pgbouncer support
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings for pgbouncer (transaction mode)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(2)

	log.Println("Database connected successfully")
	return db, nil
}
