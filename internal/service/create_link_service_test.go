package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Foga2H/ya-go-url-shortener/internal/logger"
	"github.com/Foga2H/ya-go-url-shortener/internal/mocks"
	"github.com/Foga2H/ya-go-url-shortener/internal/storage/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateLinkService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		storage := mocks.NewStorageRepo(t)
		storage.EXPECT().
			Set(mock.Anything, "user-1", mock.AnythingOfType("string"), "http://example.com").
			Return("abc123", nil)

		svc := NewCreateLinkService(storage, "http://localhost:8080", logger.NewLogger())

		result, err := svc.Create(context.Background(), "user-1", "http://example.com")
		require.NoError(t, err)
		assert.False(t, result.IsConflict)
		assert.Contains(t, result.ShortURL, "http://localhost:8080/")
	})

	t.Run("invalid url", func(t *testing.T) {
		storage := mocks.NewStorageRepo(t)
		svc := NewCreateLinkService(storage, "http://localhost:8080", logger.NewLogger())

		_, err := svc.Create(context.Background(), "user-1", "not-a-url")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidURL)
	})

	t.Run("conflict", func(t *testing.T) {
		storage := mocks.NewStorageRepo(t)
		storage.EXPECT().
			Set(mock.Anything, "user-1", mock.AnythingOfType("string"), "http://example.com").
			Return("fixed", db.ErrOriginalURLConflict)

		svc := NewCreateLinkService(storage, "http://localhost:8080", logger.NewLogger())

		result, err := svc.Create(context.Background(), "user-1", "http://example.com")
		require.NoError(t, err)
		assert.True(t, result.IsConflict)
		assert.Equal(t, "http://localhost:8080/fixed", result.ShortURL)
	})

	t.Run("save error", func(t *testing.T) {
		storage := mocks.NewStorageRepo(t)
		storage.EXPECT().
			Set(mock.Anything, "user-1", mock.AnythingOfType("string"), "http://example.com").
			Return("", errors.New("db down"))

		svc := NewCreateLinkService(storage, "http://localhost:8080", logger.NewLogger())

		_, err := svc.Create(context.Background(), "user-1", "http://example.com")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrSaveLink)
	})
}
