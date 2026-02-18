package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Foga2H/ya-go-url-shortener/internal/logger"
	"github.com/Foga2H/ya-go-url-shortener/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestShortenBatchService_ShortenError(t *testing.T) {
	storage := mocks.NewStorageRepo(t)
	storage.EXPECT().
		Set(mock.Anything, "user-1", mock.AnythingOfType("string"), "http://example.com").
		Return("", errors.New("db error"))
	shortenSvc := NewShortenService(storage, "http://localhost:8080", logger.NewLogger())
	batchSvc := NewShortenBatchService(shortenSvc)

	_, err := batchSvc.Shorten(context.Background(), "user-1", []BatchItem{{CorrelationID: "1", OriginalURL: "http://example.com"}})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrSaveShortLink)
}
