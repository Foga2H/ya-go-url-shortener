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

var ErrInvalidURL = errors.New("invalid url")
var ErrGenerateShortLink = errors.New("generate short link")
var ErrSaveLink = errors.New("save link")

type CreateLinkService struct {
	storage   repository.StorageRepo
	prefixURL string
	logger    *logger.Logger
}

type CreateLinkResult struct {
	ShortURL   string
	IsConflict bool
}

func NewCreateLinkService(storage repository.StorageRepo, prefixURL string, logger *logger.Logger) *CreateLinkService {
	return &CreateLinkService{
		storage:   storage,
		prefixURL: prefixURL,
		logger:    logger,
	}
}

func (s *CreateLinkService) Create(ctx context.Context, userID, originalURL string) (CreateLinkResult, error) {
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		s.logger.Errorf("Invalid URL: %s", originalURL)
		return CreateLinkResult{}, ErrInvalidURL
	}

	randomString, err := utils.GenerateRandomStringURLSafe(6)
	if err != nil {
		s.logger.Errorf("Failed to generate short link: %v", err)
		return CreateLinkResult{}, ErrGenerateShortLink
	}

	storedKey, err := s.storage.Set(ctx, userID, randomString, originalURL)
	if err != nil {
		if errors.Is(err, db.ErrOriginalURLConflict) {
			shortURL := s.joinShortURL(storedKey)
			s.logger.Infof("URL conflict for %s, reusing key %s", originalURL, storedKey)
			return CreateLinkResult{
				ShortURL:   shortURL,
				IsConflict: true,
			}, nil
		}
		s.logger.Errorf("Failed to save link for %s: %v", originalURL, err)
		return CreateLinkResult{}, ErrSaveLink
	}

	shortURL := s.joinShortURL(storedKey)
	s.logger.Infof("Generated link %s for %s", shortURL, originalURL)

	return CreateLinkResult{
		ShortURL: shortURL,
	}, nil
}

func (s *CreateLinkService) joinShortURL(storedKey string) string {
	shortURL, err := url.JoinPath(s.prefixURL, storedKey)
	if err != nil {
		s.logger.Errorf("Failed to join short URL parts: prefix=%s key=%s err=%v", s.prefixURL, storedKey, err)
		return s.prefixURL + "/" + storedKey
	}

	return shortURL
}
