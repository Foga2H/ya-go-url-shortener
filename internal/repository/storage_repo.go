package repository

import "context"

type UserLink struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
	IsDeleted   bool   `json:"is_deleted"`
}

type StorageRepo interface {
	Set(ctx context.Context, userID, key, value string) (string, error)
	Get(ctx context.Context, key string) (UserLink, bool)
	GetByUserID(ctx context.Context, userID string) ([]UserLink, error)
	BatchDelete(ctx context.Context, userID string, links []string) error
}
