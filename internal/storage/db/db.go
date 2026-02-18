package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

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
	IsDeleted   bool   `sql:"is_deleted"`
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

func (s *Storage) Get(ctx context.Context, key string) (repository.UserLink, bool) {
	row := s.db.QueryRowContext(ctx, "SELECT original_url, is_deleted FROM links WHERE short_url = $1", key)

	var item StorageItem
	err := row.Scan(&item.OriginalURL, &item.IsDeleted)
	if err != nil {
		return repository.UserLink{}, false
	}

	return repository.UserLink{
		ShortURL:    key,
		OriginalURL: item.OriginalURL,
		IsDeleted:   item.IsDeleted,
	}, true
}

func (s *Storage) GetByUserID(ctx context.Context, userID string) ([]repository.UserLink, error) {
	rows, err := s.db.QueryContext(
		ctx,
		"SELECT short_url, original_url FROM links WHERE user_id = $1 AND is_deleted = FALSE",
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

func (s *Storage) BatchDelete(ctx context.Context, userID string, links []string) error {
	if len(links) == 0 {
		return nil
	}

	args := make([]any, 0, len(links)+1)
	args = append(args, userID)

	placeholders := make([]string, 0, len(links))
	for i, shortURL := range links {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+2))
		args = append(args, shortURL)
	}

	query := fmt.Sprintf(
		"UPDATE links SET is_deleted = TRUE WHERE user_id = $1 AND short_url IN (%s)",
		strings.Join(placeholders, ", "),
	)

	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}
