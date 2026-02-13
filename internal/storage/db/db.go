package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
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

func (s *Storage) Set(ctx context.Context, userID, key, value string) (string, error) {
	var shortURL string
	err := s.db.QueryRowContext(
		ctx,
		`INSERT INTO links (uuid, user_id, short_url, original_url)
		 VALUES ($1, $2, $3, $4)
		 RETURNING short_url`,
		uuid.NewString(),
		userID,
		key,
		value,
	).Scan(&shortURL)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			err = s.db.QueryRowContext(
				ctx,
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

func (s *Storage) Get(ctx context.Context, key string) (string, bool) {
	row := s.db.QueryRowContext(ctx, "SELECT original_url FROM links WHERE short_url = $1", key)

	var item StorageItem
	err := row.Scan(&item.OriginalURL)
	if err != nil {
		return "", false
	}

	return item.OriginalURL, true
}

func (s *Storage) GetByUserID(ctx context.Context, userID string) ([]repository.UserLink, error) {
	rows, err := s.db.QueryContext(
		ctx,
		"SELECT short_url, original_url FROM links WHERE user_id = $1",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]repository.UserLink, 0)
	for rows.Next() {
		var item repository.UserLink
		if err := rows.Scan(&item.ShortURL, &item.OriginalURL); err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
