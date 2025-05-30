package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"

	"github.com/DaniilKalts/url-shortener/internal/storage"
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
			id INT PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
		    url TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_alias ON urls(alias);
	`,
	)
	if err != nil {
		return nil, fmt.Errorf("%s - %w", operation, err)
	}
	defer statement.Close()

	if _, err := statement.Exec(); err != nil {
		return nil, fmt.Errorf("%s - %w", operation, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(alias string, url string) (int, error) {
	const operation = "storage.sqlite.SaveURL"

	statement, err := s.db.Prepare(`INSERT INTO urls(alias, url) VALUES(?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("%s - %w", operation, err)
	}
	defer statement.Close()

	result, err := statement.Exec(alias, url)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s - %w", operation, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s - %w", operation, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s - %w", operation, err)
	}

	return int(id), nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const operation = "storage.sqlite.GetURL"

	statement, err := s.db.Prepare(`SELECT url FROM urls WHERE alias = ?`)
	if err != nil {
		return "", fmt.Errorf("%s - %w", operation, err)
	}
	defer statement.Close()

	var result string
	if err := statement.QueryRow(alias).Scan(&result); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s - %w", operation, err)
	}

	return result, nil
}
