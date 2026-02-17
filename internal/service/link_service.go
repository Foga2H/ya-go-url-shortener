package service

import (
	"context"
	"errors"

	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
)

var ErrLinkNotFound = errors.New("link not found")
var ErrLinkIsDeleted = errors.New("link is deleted")

type LinkService struct {
	storage repository.StorageRepo
}

func NewLinkService(storage repository.StorageRepo) *LinkService {
	return &LinkService{storage: storage}
}

func (s *LinkService) Resolve(ctx context.Context, id string) (string, error) {
	link, ok := s.storage.Get(ctx, id)
	if !ok {
		return "", ErrLinkNotFound
	}

	if link.IsDeleted {
		return "", ErrLinkIsDeleted
	}

	return link.OriginalURL, nil
}
