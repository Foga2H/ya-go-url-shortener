package service

import (
	"context"
	"strings"

	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
)

type UserURLsService struct {
	storage repository.StorageRepo
	prefix  string
}

type UserURL struct {
	ShortURL    string
	OriginalURL string
}

func NewUserURLsService(storage repository.StorageRepo, prefixURL string) *UserURLsService {
	return &UserURLsService{storage: storage, prefix: strings.TrimRight(prefixURL, "/") + "/"}
}

func (s *UserURLsService) List(ctx context.Context, userID string) ([]UserURL, error) {
	results, err := s.storage.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := make([]UserURL, 0, len(results))
	for _, item := range results {
		response = append(response, UserURL{
			ShortURL:    s.prefix + strings.TrimLeft(item.ShortURL, "/"),
			OriginalURL: item.OriginalURL,
		})
	}

	return response, nil
}
