package service

import (
	"context"
	"testing"

	"github.com/Foga2H/ya-go-url-shortener/internal/logger"
	"github.com/Foga2H/ya-go-url-shortener/internal/mocks"
	"github.com/Foga2H/ya-go-url-shortener/internal/storage/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestShortenService_ShortenConflict(t *testing.T) {
	storage := mocks.NewStorageRepo(t)
	storage.EXPECT().
		Set(mock.Anything, "user-1", mock.AnythingOfType("string"), "http://example.com").
		Return("known", db.ErrOriginalURLConflict)

	svc := NewShortenService(storage, "http://localhost:8080", logger.NewLogger())

	result, err := svc.Shorten(context.Background(), "user-1", "http://example.com")
	require.NoError(t, err)
	assert.True(t, result.IsConflict)
	assert.Equal(t, "http://localhost:8080/known", result.ShortURL)
}
