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

func TestUserURLsService_List(t *testing.T) {
	storage := mocks.NewStorageRepo(t)
	storage.EXPECT().
		GetByUserID(mock.Anything, "user-1").
		Return([]repository.UserLink{{ShortURL: "/abc", OriginalURL: "http://example.com"}}, nil)
	svc := NewUserURLsService(storage, "http://localhost:8080/", logger.NewLogger())

	items, err := svc.List(context.Background(), "user-1")
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "http://localhost:8080/abc", items[0].ShortURL)
	assert.Equal(t, "http://example.com", items[0].OriginalURL)
}
