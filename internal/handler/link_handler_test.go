package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Foga2H/ya-go-url-shortener/internal/config"
	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
	storage "github.com/Foga2H/ya-go-url-shortener/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinkHandler_ServeHTTP(t *testing.T) {
	memStorage := storage.NewMemStorage()
	_, err := memStorage.Set(`test`, `http://yandex.ru`)
	require.NoError(t, err)

	type fields struct {
		Storage repository.StorageRepo
		config  *config.Config
	}
	type args struct {
		url  string
		body io.Reader
	}
	type want struct {
		code        int
		response    string
		contentType string
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
				Storage: memStorage,
			},
			args: args{
				url:  "/test",
				body: nil,
			},
			want: want{
				code:        http.StatusTemporaryRedirect,
				response:    "http://yandex.ru",
				contentType: "text/plain",
			},
		},
		{
			name: "not found",
			fields: fields{
				Storage: memStorage,
			},
			args: args{
				url:  "/awdwadaw",
				body: nil,
			},
			want: want{
				code:        http.StatusNotFound,
				response:    "Link not found",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &LinkHandler{
				Storage: tt.fields.Storage,
			}
			request := httptest.NewRequest(http.MethodGet, tt.args.url, tt.args.body)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			h.ServeHTTP(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			strData := string(resBody)

			trimmedData := strings.TrimSuffix(strData, "\n")
			trimmedData = strings.TrimSuffix(trimmedData, "\r")

			require.NoError(t, err)
			assert.Equal(t, tt.want.response, trimmedData)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
