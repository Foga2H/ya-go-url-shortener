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

type shortenStorageMock struct {
	set func(ctx context.Context, userID, key, value string) (string, error)
}

func (m *shortenStorageMock) Set(ctx context.Context, userID, key, value string) (string, error) {
	return m.set(ctx, userID, key, value)
}

func (m *shortenStorageMock) Get(_ context.Context, _ string) (string, bool) {
	return "", false
}

func (m *shortenStorageMock) GetByUserID(_ context.Context, _ string) ([]repository.UserLink, error) {
	return nil, nil
}

func TestShortenService_ShortenConflict(t *testing.T) {
	svc := NewShortenService(&shortenStorageMock{
		set: func(_ context.Context, _ string, _ string, _ string) (string, error) {
			return "known", db.ErrOriginalURLConflict
		},
	}, "http://localhost:8080")

	result, err := svc.Shorten(context.Background(), "user-1", "http://example.com")
	require.NoError(t, err)
	assert.True(t, result.IsConflict)
	assert.Equal(t, "http://localhost:8080/known", result.ShortURL)
}

func TestShortenBatchService_ShortenError(t *testing.T) {
	shortenSvc := NewShortenService(&shortenStorageMock{
		set: func(_ context.Context, _ string, _ string, _ string) (string, error) {
			return "", errors.New("db error")
		},
	}, "http://localhost:8080")
	batchSvc := NewShortenBatchService(shortenSvc)

	_, err := batchSvc.Shorten(context.Background(), "user-1", []BatchItem{{CorrelationID: "1", OriginalURL: "http://example.com"}})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrSaveShortLink)
}

type linkStorageMock struct {
	link string
	ok   bool
}

func (m *linkStorageMock) Set(_ context.Context, _ string, key, _ string) (string, error) {
	return key, nil
}

func (m *linkStorageMock) Get(_ context.Context, _ string) (string, bool) {
	return m.link, m.ok
}

func (m *linkStorageMock) GetByUserID(_ context.Context, _ string) ([]repository.UserLink, error) {
	return nil, nil
}

func TestLinkService_ResolveNotFound(t *testing.T) {
	svc := NewLinkService(&linkStorageMock{ok: false})

	_, err := svc.Resolve(context.Background(), "abc")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrLinkNotFound)
}

type userURLsStorageMock struct {
	items []repository.UserLink
	err   error
}

func (m *userURLsStorageMock) Set(_ context.Context, _ string, key, _ string) (string, error) {
	return key, nil
}

func (m *userURLsStorageMock) Get(_ context.Context, _ string) (string, bool) {
	return "", false
}

func (m *userURLsStorageMock) GetByUserID(_ context.Context, _ string) ([]repository.UserLink, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.items, nil
}

func TestUserURLsService_List(t *testing.T) {
	svc := NewUserURLsService(&userURLsStorageMock{items: []repository.UserLink{{ShortURL: "/abc", OriginalURL: "http://example.com"}}}, "http://localhost:8080/")

	items, err := svc.List(context.Background(), "user-1")
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "http://localhost:8080/abc", items[0].ShortURL)
	assert.Equal(t, "http://example.com", items[0].OriginalURL)
}

type pingerMock struct {
	err error
}

func (m *pingerMock) PingContext(_ context.Context) error {
	return m.err
}

func TestPingService_PingError(t *testing.T) {
	svc := NewPingService(&pingerMock{err: errors.New("unavailable")})

	err := svc.Ping(context.Background())
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrPingDatabase)
}
