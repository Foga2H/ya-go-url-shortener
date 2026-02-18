package handler

import (
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Foga2H/ya-go-url-shortener/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPingHandler(t *testing.T) {
	dbConn, err := sql.Open("pgx", "postgres://127.0.0.1:5432/test?sslmode=disable")
	require.NoError(t, err)
	defer dbConn.Close()

	h := NewPingHandler(dbConn, logger.NewLogger())

	require.NotNil(t, h)
}

func TestPingHandler_ServeHTTP_DatabaseUnavailable(t *testing.T) {
	dbConn, err := sql.Open("pgx", "postgres://127.0.0.1:1/test?sslmode=disable")
	require.NoError(t, err)
	defer dbConn.Close()

	h := NewPingHandler(dbConn, logger.NewLogger())

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Contains(t, res.Header.Get("Content-Type"), "text/plain")

	bodyBytes, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Contains(t, strings.TrimSpace(string(bodyBytes)), http.StatusText(http.StatusInternalServerError))
}
