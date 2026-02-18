package middleware

import (
	"context"
	"net/http"

	"github.com/Foga2H/ya-go-url-shortener/internal/auth"
	"github.com/google/uuid"
)

const userCookieName = "user"

type UserCookie struct {
	token *auth.Token
}

func NewUserCookie(token *auth.Token) *UserCookie {
	return &UserCookie{
		token: token,
	}
}

type contextKey string

const userIDKey contextKey = "userID"

func ContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(userIDKey).(string)
	return v, ok
}

func (uc *UserCookie) Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := ""

			userCookie, err := r.Cookie(userCookieName)
			if err == nil {
				userID, err = uc.token.ParseUserID(userCookie.Value)
			}

			if err != nil || userID == "" {
				userID = uuid.NewString()
				token, tokenErr := uc.token.GenerateToken(userID)
				if tokenErr != nil {
					http.Error(w, "failed to generate user token", http.StatusInternalServerError)
					return
				}

				http.SetCookie(w, &http.Cookie{
					Name:     userCookieName,
					Value:    token,
					Path:     "/",
					HttpOnly: true,
					MaxAge:   int(auth.TokenExp.Seconds()),
					SameSite: http.SameSiteLaxMode,
				})
			}

			next.ServeHTTP(w, r.WithContext(ContextWithUserID(r.Context(), userID)))
		})
	}
}
