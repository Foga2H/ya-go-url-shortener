package service

import (
	"context"
	"errors"
	"net/url"

	"github.com/Foga2H/ya-go-url-shortener/internal/logger"
	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
	"github.com/Foga2H/ya-go-url-shortener/internal/storage/db"
	"github.com/Foga2H/ya-go-url-shortener/pkg/utils"
)

var ErrSaveShortLink = errors.New("save short link")

type ShortenService struct {
	storage   repository.StorageRepo
	prefixURL string
	logger    *logger.Logger
}

type ShortenResult struct {
	ShortURL   string
	IsConflict bool
}

func NewShortenService(storage repository.StorageRepo, prefixURL string, logger *logger.Logger) *ShortenService {
	return &ShortenService{storage: storage, prefixURL: prefixURL, logger: logger}
}

func (s *ShortenService) Shorten(ctx context.Context, userID, originalURL string) (ShortenResult, error) {
	randomString, err := utils.GenerateRandomStringURLSafe(6)
	if err != nil {
		s.logger.Errorf("Failed to generate short link: %v", err)
		return ShortenResult{}, ErrGenerateShortLink
	}

	storedKey, err := s.storage.Set(ctx, userID, randomString, originalURL)
	if err != nil {
		if errors.Is(err, db.ErrOriginalURLConflict) {
			shortURL := s.joinShortURL(storedKey)
			s.logger.Infof("URL conflict for %s, reusing key %s", originalURL, storedKey)
			return ShortenResult{ShortURL: shortURL, IsConflict: true}, nil
		}
		s.logger.Errorf("Failed to save short link for %s: %v", originalURL, err)
		return ShortenResult{}, ErrSaveShortLink
	}

	shortURL := s.joinShortURL(storedKey)
	s.logger.Infof("Generated link %s for %s", shortURL, originalURL)

	return ShortenResult{ShortURL: shortURL}, nil
}

func (s *ShortenService) joinShortURL(storedKey string) string {
	shortURL, err := url.JoinPath(s.prefixURL, storedKey)
	if err != nil {
		s.logger.Errorf("Failed to join short URL parts: prefix=%s key=%s err=%v", s.prefixURL, storedKey, err)
		return s.prefixURL + "/" + storedKey
	}

	return shortURL
}
