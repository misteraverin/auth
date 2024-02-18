package token

import (
	"BasicAuth/internal/errdomain"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

const secret_key string = "your-256-bit-secret" //"some_secret_key"

type Encrypted struct {
	Token string `json:token`
}

type JWTToken struct {
	Login     string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

func New(login string, duration time.Duration) (*JWTToken, error) {
	return NewWithStartTime(login, time.Now(), duration)
}

func NewWithStartTime(login string, startTime time.Time, duration time.Duration) (*JWTToken, error) {
	token := JWTToken{
		Login:     login,
		IssuedAt:  startTime,
		ExpiresAt: startTime.Add(duration),
	}

	return &token, nil
}

func (t *JWTToken) Valid() error {
	if time.Now().After(t.ExpiresAt) {
		return errdomain.ErrTokenExpired
	}

	return nil
}

func (t *JWTToken) Encrypt() (*Encrypted, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t)

	signingKey := []byte(secret_key)
	tokenString, err := token.SignedString(signingKey)

	return &Encrypted{Token: tokenString}, err
}

func ParseJWT(tokenString string) (*JWTToken, error) {
	keyfunc := func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errdomain.ErrTokenInvalid
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret_key), nil
	}

	jwtToken := JWTToken{}
	_, err := jwt.ParseWithClaims(tokenString, &jwtToken, keyfunc)

	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, errdomain.ErrTokenExpired) {
			return nil, errdomain.ErrTokenExpired
		}

		return nil, errdomain.ErrTokenInvalid
	}

	return &jwtToken, nil
}

func (t *JWTToken) Update(duration time.Duration) (*Encrypted, error) {
	t, err := New(t.Login, duration)
	if err != nil {
		return nil, err
	}

	encryptedToken, err := t.Encrypt()
	if err != nil {
		return nil, err
	}

	return encryptedToken, nil
}
