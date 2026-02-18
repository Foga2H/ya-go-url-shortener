package service

import (
	"context"
	"errors"

	"github.com/Foga2H/ya-go-url-shortener/internal/logger"
	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
)

var ErrLinkNotFound = errors.New("link not found")
var ErrLinkIsDeleted = errors.New("link is deleted")

type LinkService struct {
	storage repository.StorageRepo
	logger  *logger.Logger
}

func NewLinkService(storage repository.StorageRepo, logger *logger.Logger) *LinkService {
	return &LinkService{storage: storage, logger: logger}
}

func (s *LinkService) Resolve(ctx context.Context, id string) (string, error) {
	link, ok := s.storage.Get(ctx, id)
	if !ok {
		s.logger.Warnf("Link not found: id=%s", id)
		return "", ErrLinkNotFound
	}

	if link.IsDeleted {
		s.logger.Warnf("Link is deleted: id=%s", id)
		return "", ErrLinkIsDeleted
	}

	return link.OriginalURL, nil
}
