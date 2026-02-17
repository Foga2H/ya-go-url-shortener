package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Foga2H/ya-go-url-shortener/internal/config/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPingHandler(t *testing.T) {
	cfg := db.NewConfigFrom("postgres://127.0.0.1:5432/test?sslmode=disable")

	h := NewPingHandler(cfg)

	require.NotNil(t, h)
	assert.Same(t, cfg, h.Config)
}

func TestPingHandler_ServeHTTP_DatabaseUnavailable(t *testing.T) {
	cfg := db.NewConfigFrom("postgres://127.0.0.1:1/test?sslmode=disable")
	h := NewPingHandler(cfg)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Contains(t, res.Header.Get("Content-Type"), "text/plain")

	bodyBytes, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Contains(t, strings.TrimSpace(string(bodyBytes)), "Error when trying to connect to database")
}
