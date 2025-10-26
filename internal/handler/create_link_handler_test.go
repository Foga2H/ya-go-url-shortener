package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/Foga2H/ya-go-url-shortener/internal/config"
	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
	storage "github.com/Foga2H/ya-go-url-shortener/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateLinkHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		Storage repository.StorageRepo
		*config.Config
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
				Config:  config.NewConfig(),
			},
			args: args{
				bodyString: `http://yahoo.com`,
			},
			want: want{
				code:        http.StatusCreated,
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateLinkHandler{
				Storage: tt.fields.Storage,
				config:  tt.fields.Config,
			}
			request := httptest.NewRequest(http.MethodPost, `/`, strings.NewReader(tt.args.bodyString))
			// создаём новый Recorder
			w := httptest.NewRecorder()
			h.ServeHTTP(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()

			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))

			bodyBytes, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			body := strings.TrimSpace(string(bodyBytes))
			t.Logf("response: %s", body)

			// Проверяем, что строка соответствует формату http://localhost:8080/<идентификатор>
			matched, err := regexp.MatchString(`^`+tt.fields.Config.PrefixUrl+`/[A-Za-z0-9_-]+$`, body)
			require.NoError(t, err)
			assert.True(t, matched, "Ожидался корректный короткий URL, получено: %s", body)
		})
	}
}
