package token

import (
	"auth/internal/errdomain"
	"time"

	"github.com/golang-jwt/jwt"
)

const secret_key string = "your-256-bit-secret" //"some_secret_key"

type Encrypted struct {
	Token string `json:token`
}

type JWT struct {
	Login     string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

func (t *JWT) Valid() error {
	if time.Now().After(t.ExpiresAt) {
		return errdomain.ErrTokenExpired
	}

	return nil
}

func (t *JWT) Encrypt() (*Encrypted, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t)

	signingKey := []byte(secret_key)
	tokenString, err := token.SignedString(signingKey)

	return &Encrypted{Token: tokenString}, err
}
