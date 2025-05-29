package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(storagePath string) (*Storage, error) {
	const operation = "storage.sqlite.NewStorage"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s - %w", operation, err)
	}

	statement, err := db.Prepare(
		`
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
		    url TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_alias ON urls(alias);
	`,
	)
	if err != nil {
		return nil, fmt.Errorf("%s - %w", operation, err)
	}

	if _, err := statement.Exec(); err != nil {
		return nil, fmt.Errorf("%s - %w", operation, err)
	}

	return &Storage{db: db}, nil
}
