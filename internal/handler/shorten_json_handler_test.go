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
	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortenJSONHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		Storage repository.StorageRepo
		config  *config.Config
	}
	type want struct {
		code        int
		response    string
		contentType string
	}
	type args struct {
		bodyString string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ShortenJSONHandler{
				Storage: tt.fields.Storage,
				config:  tt.fields.config,
			}

			request := httptest.NewRequest(http.MethodPost, `/api/shorten`, strings.NewReader(tt.args.bodyString))

			w := httptest.NewRecorder()
			h.ServeHTTP(w, request)

			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)

			defer res.Body.Close()

			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))

			bodyBytes, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			var resp Response
			err = json.Unmarshal(bodyBytes, &resp)
			t.Logf("response: %s", string(bodyBytes))
			require.NoError(t, err)

			// Проверяем, что строка соответствует формату http://localhost:8080/<идентификатор>
			matched, err := regexp.MatchString(`^`+tt.fields.config.PrefixURL+`/[A-Za-z0-9_-]+$`, resp.Result)
			require.NoError(t, err)
			assert.True(t, matched, "Ожидался корректный короткий URL, получено: %s", resp.Result)
		})
	}
}
