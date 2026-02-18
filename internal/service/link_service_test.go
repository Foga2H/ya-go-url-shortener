package service

import (
	"context"
	"testing"

	"github.com/Foga2H/ya-go-url-shortener/internal/logger"
	"github.com/Foga2H/ya-go-url-shortener/internal/mocks"
	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLinkService_ResolveNotFound(t *testing.T) {
	storage := mocks.NewStorageRepo(t)
	storage.EXPECT().
		Get(mock.Anything, "abc").
		Return(repository.UserLink{}, false)
	svc := NewLinkService(storage, logger.NewLogger())

	_, err := svc.Resolve(context.Background(), "abc")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrLinkNotFound)
}

func TestLinkService_ResolveDeleted(t *testing.T) {
	storage := mocks.NewStorageRepo(t)
	storage.EXPECT().
		Get(mock.Anything, "abc").
		Return(repository.UserLink{
			OriginalURL: "http://example.com",
			IsDeleted:   true,
		}, true)
	svc := NewLinkService(storage, logger.NewLogger())

	_, err := svc.Resolve(context.Background(), "abc")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrLinkIsDeleted)
}
