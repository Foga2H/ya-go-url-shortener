package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/Foga2H/ya-go-url-shortener/internal/config"
	"github.com/Foga2H/ya-go-url-shortener/internal/storage/db"
	storage "github.com/Foga2H/ya-go-url-shortener/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortenBatchJSONHandler_ServeHTTP_Success(t *testing.T) {
	memStorage := storage.NewMemStorage()
	cfg := config.NewConfigFrom("localhost:8080", "http://localhost:8080")
	h := NewShortenBatchJSONHandler(memStorage, cfg)

	body := `[
		{"correlation_id":"id-1","original_url":"http://yandex.ru"},
		{"correlation_id":"id-2","original_url":"http://google.com"}
	]`

	req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(body))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	bodyBytes, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var resp []BatchResult
	require.NoError(t, json.Unmarshal(bodyBytes, &resp))
	require.Len(t, resp, 2)

	assert.Equal(t, "id-1", resp[0].CorrelationID)
	assert.Equal(t, "id-2", resp[1].CorrelationID)

	for i, item := range resp {
		matched, err := regexp.MatchString(`^`+cfg.PrefixURL+`/[A-Za-z0-9_-]+$`, item.ShortURL)
		require.NoError(t, err)
		assert.True(t, matched, "Ожидался корректный короткий URL, получено: %s", item.ShortURL)

		key := strings.TrimPrefix(item.ShortURL, cfg.PrefixURL+"/")
		got, ok := memStorage.Get(key)
		require.True(t, ok)
		if i == 0 {
			assert.Equal(t, "http://yandex.ru", got)
			continue
		}
		assert.Equal(t, "http://google.com", got)
	}
}

func TestShortenBatchJSONHandler_ServeHTTP_InvalidJSON(t *testing.T) {
	memStorage := storage.NewMemStorage()
	cfg := config.NewConfigFrom("localhost:8080", "http://localhost:8080")
	h := NewShortenBatchJSONHandler(memStorage, cfg)

	req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(`{"invalid"`))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Contains(t, res.Header.Get("Content-Type"), "text/plain")

	bodyBytes, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Contains(t, strings.TrimSpace(string(bodyBytes)), "Invalid JSON Body")
}

type conflictStorage struct {
	originalToKey map[string]string
}

func (s *conflictStorage) Set(key string, value string) (string, error) {
	if existingKey, ok := s.originalToKey[value]; ok {
		return existingKey, db.ErrOriginalURLConflict
	}
	s.originalToKey[value] = key
	return key, nil
}

func (s *conflictStorage) Get(key string) (string, bool) {
	for originalURL, storedKey := range s.originalToKey {
		if storedKey == key {
			return originalURL, true
		}
	}
	return "", false
}

func TestShortenBatchJSONHandler_ServeHTTP_ConflictInBatch(t *testing.T) {
	cfg := config.NewConfigFrom("localhost:8080", "http://localhost:8080")
	st := &conflictStorage{
		originalToKey: map[string]string{
			"http://yandex.ru": "exists1",
		},
	}
	h := NewShortenBatchJSONHandler(st, cfg)

	body := `[
		{"correlation_id":"id-1","original_url":"http://yandex.ru"},
		{"correlation_id":"id-2","original_url":"http://google.com"}
	]`

	req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(body))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	bodyBytes, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var resp []BatchResult
	require.NoError(t, json.Unmarshal(bodyBytes, &resp))
	require.Len(t, resp, 2)

	assert.Equal(t, "id-1", resp[0].CorrelationID)
	assert.Equal(t, cfg.PrefixURL+"/exists1", resp[0].ShortURL)
	assert.Equal(t, "id-2", resp[1].CorrelationID)

	assert.NotEqual(t, cfg.PrefixURL+"/exists1", resp[1].ShortURL)
}
