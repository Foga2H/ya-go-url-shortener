package repository

import "context"

type UserLink struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
	IsDeleted   bool   `json:"is_deleted"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --dir . --name StorageRepo --output ../mocks --outpkg mocks --filename storage_repo.go --with-expecter --disable-version-string --issue-845-fix
type StorageRepo interface {
	Set(ctx context.Context, userID, key, value string) (string, error)
	Get(ctx context.Context, key string) (UserLink, bool)
	GetByUserID(ctx context.Context, userID string) ([]UserLink, error)
	BatchDelete(ctx context.Context, userID string, links []string) error
}
