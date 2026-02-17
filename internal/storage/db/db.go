package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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

var ErrOriginalURLConflict = errors.New("original_url already exists")

func (s *Storage) Set(key string, value string) (string, error) {
	var shortURL string
	err := s.db.QueryRowContext(
		context.Background(),
		`INSERT INTO links (uuid, short_url, original_url)
		 VALUES ($1, $2, $3)
		 RETURNING short_url`,
		uuid.New().String(),
		key,
		value,
	).Scan(&shortURL)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			err = s.db.QueryRowContext(
				context.Background(),
				"SELECT short_url FROM links WHERE original_url = $1",
				value,
			).Scan(&shortURL)
			if err != nil {
				return "", err
			}
			return shortURL, ErrOriginalURLConflict
		}

		return "", err
	}

	return shortURL, nil
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
