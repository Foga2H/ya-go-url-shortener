package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
	"github.com/Foga2H/ya-go-url-shortener/internal/storage/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testStorage struct {
	setFunc func(ctx context.Context, userID, key, value string) (string, error)
}

func (s *testStorage) Set(ctx context.Context, userID, key, value string) (string, error) {
	return s.setFunc(ctx, userID, key, value)
}

func (s *testStorage) Get(_ context.Context, _ string) (string, bool) {
	return "", false
}

func (s *testStorage) GetByUserID(_ context.Context, _ string) ([]repository.UserLink, error) {
	return nil, nil
}

func TestCreateLinkService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := NewCreateLinkService(&testStorage{
			setFunc: func(_ context.Context, _ string, key, _ string) (string, error) {
				return key, nil
			},
		}, "http://localhost:8080")

		result, err := svc.Create(context.Background(), "user-1", "http://example.com")
		require.NoError(t, err)
		assert.False(t, result.IsConflict)
		assert.Contains(t, result.ShortURL, "http://localhost:8080/")
	})

	t.Run("invalid url", func(t *testing.T) {
		svc := NewCreateLinkService(&testStorage{
			setFunc: func(_ context.Context, _ string, key, _ string) (string, error) {
				return key, nil
			},
		}, "http://localhost:8080")

		_, err := svc.Create(context.Background(), "user-1", "not-a-url")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidURL)
	})

	t.Run("conflict", func(t *testing.T) {
		svc := NewCreateLinkService(&testStorage{
			setFunc: func(_ context.Context, _ string, _ string, _ string) (string, error) {
				return "fixed", db.ErrOriginalURLConflict
			},
		}, "http://localhost:8080")

		result, err := svc.Create(context.Background(), "user-1", "http://example.com")
		require.NoError(t, err)
		assert.True(t, result.IsConflict)
		assert.Equal(t, "http://localhost:8080/fixed", result.ShortURL)
	})

	t.Run("save error", func(t *testing.T) {
		svc := NewCreateLinkService(&testStorage{
			setFunc: func(_ context.Context, _ string, _ string, _ string) (string, error) {
				return "", errors.New("db down")
			},
		}, "http://localhost:8080")

		_, err := svc.Create(context.Background(), "user-1", "http://example.com")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrSaveLink)
	})
}
