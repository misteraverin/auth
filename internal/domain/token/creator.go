package token

import (
	"auth/internal/domain/clock"
	"auth/internal/errdomain"
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

type Creator struct {
	clock       clock.Interface
	expiredTime time.Duration
}

func NewCreator(clock clock.Interface, expiredTime time.Duration) *Creator {
	return &Creator{clock: clock, expiredTime: expiredTime}
}

func (c *Creator) NewJWT(login string) (*JWT, error) {
	token := JWT{
		Login:     login,
		IssuedAt:  c.clock.Now(),
		ExpiresAt: c.clock.Now().Add(c.expiredTime),
	}

	return &token, nil
}

func (c *Creator) ParseJWT(tokenString string) (*JWT, error) {
	keyfunc := func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errdomain.ErrTokenInvalid
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret_key), nil
	}

	jwtToken := JWT{}
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

func (c *Creator) UpdateJWT(login string) (*Encrypted, error) {
	t, err := c.NewJWT(login)
	if err != nil {
		return nil, err
	}

	encryptedToken, err := t.Encrypt()
	if err != nil {
		return nil, err
	}

	return encryptedToken, nil
}
