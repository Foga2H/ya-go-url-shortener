package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"

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
}

type CreateLinkResult struct {
	ShortURL   string
	IsConflict bool
}

func NewCreateLinkService(storage repository.StorageRepo, prefixURL string) *CreateLinkService {
	return &CreateLinkService{
		storage:   storage,
		prefixURL: prefixURL,
	}
}

func (s *CreateLinkService) Create(ctx context.Context, userID, originalURL string) (CreateLinkResult, error) {
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		return CreateLinkResult{}, ErrInvalidURL
	}

	randomString, err := utils.GenerateRandomStringURLSafe(6)
	if err != nil {
		return CreateLinkResult{}, ErrGenerateShortLink
	}

	storedKey, err := s.storage.Set(ctx, userID, randomString, originalURL)
	if err != nil {
		if errors.Is(err, db.ErrOriginalURLConflict) {
			return CreateLinkResult{
				ShortURL:   s.prefixURL + "/" + storedKey,
				IsConflict: true,
			}, nil
		}
		return CreateLinkResult{}, ErrSaveLink
	}

	fmt.Printf("Generated link %s for %s\n", s.prefixURL+"/"+storedKey, originalURL)

	return CreateLinkResult{
		ShortURL: s.prefixURL + "/" + storedKey,
	}, nil
}
