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
	"github.com/Foga2H/ya-go-url-shortener/internal/logger"
	"github.com/Foga2H/ya-go-url-shortener/internal/middleware"
	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
	storage "github.com/Foga2H/ya-go-url-shortener/internal/storage/memory"
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
		{
			name: "success",
			fields: fields{
				Storage: storage.NewMemStorage(),
				config:  config.NewConfigFrom("localhost:8080", "http://localhost:8080"),
			},
			args: args{
				bodyString: `{"url":"http://yandex.ru"}`,
			},
			want: want{
				code:        http.StatusCreated,
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewShortenJSONHandler(tt.fields.Storage, tt.fields.config, logger.NewLogger())

			request := httptest.NewRequest(http.MethodPost, `/api/shorten`, strings.NewReader(tt.args.bodyString))
			request = request.WithContext(middleware.ContextWithUserID(request.Context(), "test-user"))

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
