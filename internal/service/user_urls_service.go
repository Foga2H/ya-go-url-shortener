package service

import (
	"context"
	"net/url"

	"github.com/Foga2H/ya-go-url-shortener/internal/logger"
	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
)

type UserURLsService struct {
	storage repository.StorageRepo
	prefix  string
	logger  *logger.Logger
}

type UserURL struct {
	ShortURL    string
	OriginalURL string
}

func NewUserURLsService(storage repository.StorageRepo, prefixURL string, logger *logger.Logger) *UserURLsService {
	return &UserURLsService{storage: storage, prefix: prefixURL, logger: logger}
}

func (s *UserURLsService) List(ctx context.Context, userID string) ([]UserURL, error) {
	results, err := s.storage.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to get URLs for user %s: %v", userID, err)
		return nil, err
	}

	response := make([]UserURL, 0, len(results))
	for _, item := range results {
		shortURL, joinErr := url.JoinPath(s.prefix, item.ShortURL)
		if joinErr != nil {
			s.logger.Errorf("Failed to join user URL parts: prefix=%s key=%s err=%v", s.prefix, item.ShortURL, joinErr)
			return nil, joinErr
		}
		response = append(response, UserURL{
			ShortURL:    shortURL,
			OriginalURL: item.OriginalURL,
		})
	}

	return response, nil
}
