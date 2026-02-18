package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const TokenExp = 3 * time.Hour

type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"uid"`
}

type Token struct {
	secret []byte
}

func NewToken(secret string) *Token {
	return &Token{
		secret: []byte(secret),
	}
}

func (t *Token) GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userID,
	})

	return token.SignedString(t.secret)
}

func (t *Token) ParseUserID(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return t.secret, nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	if claims.UserID == "" {
		return "", errors.New("empty user id")
	}

	return claims.UserID, nil
}
