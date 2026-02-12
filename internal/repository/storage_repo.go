package repository

import "context"

type UserLink struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type StorageRepo interface {
	Set(ctx context.Context, userID, key, value string) (string, error)
	Get(ctx context.Context, key string) (string, bool)
	GetByUserID(ctx context.Context, userID string) ([]UserLink, error)
}
