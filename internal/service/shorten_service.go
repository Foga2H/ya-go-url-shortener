package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
	"github.com/Foga2H/ya-go-url-shortener/internal/storage/db"
	"github.com/Foga2H/ya-go-url-shortener/pkg/utils"
)

var ErrSaveShortLink = errors.New("save short link")

type ShortenService struct {
	storage   repository.StorageRepo
	prefixURL string
}

type ShortenResult struct {
	ShortURL   string
	IsConflict bool
}

func NewShortenService(storage repository.StorageRepo, prefixURL string) *ShortenService {
	return &ShortenService{storage: storage, prefixURL: prefixURL}
}

func (s *ShortenService) Shorten(ctx context.Context, userID, originalURL string) (ShortenResult, error) {
	randomString, err := utils.GenerateRandomStringURLSafe(6)
	if err != nil {
		return ShortenResult{}, ErrGenerateShortLink
	}

	storedKey, err := s.storage.Set(ctx, userID, randomString, originalURL)
	if err != nil {
		if errors.Is(err, db.ErrOriginalURLConflict) {
			return ShortenResult{ShortURL: s.prefixURL + "/" + storedKey, IsConflict: true}, nil
		}
		return ShortenResult{}, ErrSaveShortLink
	}

	fmt.Printf("Generated link %s for %s\n", s.prefixURL+"/"+storedKey, originalURL)

	return ShortenResult{ShortURL: s.prefixURL + "/" + storedKey}, nil
}
