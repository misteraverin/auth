package handlers_test

import (
	"BasicAuth/internal/httpserver/handlers"
	"BasicAuth/internal/httpserver/token"
	"BasicAuth/internal/repository"
	"BasicAuth/pkg/emptylogger"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLoginHandler(t *testing.T) {
	cases := []struct {
		name         string
		login        string
		password     string
		useBasicAuth bool
	}{
		{
			name:         "default case",
			login:        "login",
			password:     "password",
			useBasicAuth: true,
		},

		{
			name:         "without basic auth",
			login:        "login",
			password:     "password",
			useBasicAuth: false,
		},
	}

	mapDB, err := repository.NewMapDB()
	assert.NoError(t, err)

	log := emptylogger.New()
	startTime := time.Now()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "https://localhost:8000/login", nil)

			if tc.useBasicAuth {
				req.SetBasicAuth(tc.login, tc.password)
			}

			loginHandler := handlers.Login(log, mapDB)
			w := httptest.NewRecorder()
			loginHandler.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			statusCode := resp.StatusCode

			if tc.useBasicAuth {
				assert.Equal(t, http.StatusOK, statusCode)
			} else {
				assert.Equal(t, http.StatusUnauthorized, statusCode)
			}

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			encryptedToken := token.Encrypted{}
			err = json.Unmarshal(body, &encryptedToken)

			if tc.useBasicAuth {
				assert.NoError(t, err)

				tok, err := token.NewWithStartTime(tc.login, startTime, handlers.TokenExpiredTime)
				assert.NoError(t, err)
				etok, err := tok.Encrypt()
				assert.NoError(t, err)
				assert.Equal(t, etok.Token, encryptedToken.Token)
			} else {
				assert.Error(t, err)
				assert.Equal(t, "", encryptedToken.Token)
			}
		})
	}
}
