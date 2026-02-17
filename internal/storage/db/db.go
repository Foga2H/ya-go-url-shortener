package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

type StorageItem struct {
	OriginalURL string `sql:"original_url"`
}

func (s *Storage) Set(key string, value string) error {
	_, err := s.db.ExecContext(context.Background(), "INSERT INTO links (uuid, short_url, original_url) VALUES ($1, $2, $3)", uuid.New().String(), key, value)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Get(key string) (string, bool) {
	row := s.db.QueryRowContext(context.Background(), "SELECT original_url FROM links WHERE short_url = $1", key)

	var item StorageItem
	err := row.Scan(&item.OriginalURL)
	if err != nil {
		return "", false
	}

	return item.OriginalURL, true
}
