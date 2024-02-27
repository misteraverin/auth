package handlers_test

import (
	"auth/internal/domain/token"
	"auth/internal/domain/user/service"
	"auth/internal/httpserver/handlers"
	"auth/internal/repository/inmemory"
	"auth/pkg/logger"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type testTime struct {
	time time.Time
}

func (t *testTime) Now() time.Time {
	return t.time
}

func TestLoginSuccess(t *testing.T) {
	cases := []struct {
		name     string
		login    string
		password string
	}{
		{
			name:     "default case 1",
			login:    "login 1",
			password: "password 1",
		},

		{
			name:     "default case 2",
			login:    "login 2",
			password: "password 2",
		},
	}

	log := logger.NewEmpty()

	db, err := inmemory.NewMapDB()
	require.NoError(t, err)

	tt := testTime{time: time.Now()}
	us, err := service.NewUserService(db, &tt)

	h := handlers.New(log, us, us)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "https://localhost:8000/login", nil)
			req.SetBasicAuth(tc.login, tc.password)

			loginHandler := h.Login()
			w := httptest.NewRecorder()
			loginHandler.ServeHTTP(w, req)

			resp := w.Result()

			statusCode := resp.StatusCode
			require.Equal(t, http.StatusOK, statusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			encryptedToken := token.Encrypted{}
			err = json.Unmarshal(body, &encryptedToken)

			require.NoError(t, err)

			creator := token.NewCreator(&tt, service.TokenExpiredTime)
			tok, err := creator.NewJWT(tc.login)
			require.NoError(t, err)
			etok, err := tok.Encrypt()
			require.NoError(t, err)
			require.Equal(t, etok.Token, encryptedToken.Token)

			err = resp.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestLoginWrong(t *testing.T) {
	cases := []struct {
		name         string
		login        string
		password     string
		useBasicAuth bool
	}{
		{
			name:     "without basic auth",
			login:    "login",
			password: "password",
		},
	}

	log := logger.NewEmpty()

	db, err := inmemory.NewMapDB()
	require.NoError(t, err)

	tt := testTime{time: time.Now()}
	us, err := service.NewUserService(db, &tt)

	h := handlers.New(log, us, us)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "https://localhost:8000/login", nil)

			loginHandler := h.Login()
			w := httptest.NewRecorder()
			loginHandler.ServeHTTP(w, req)

			resp := w.Result()

			statusCode := resp.StatusCode
			require.Equal(t, http.StatusUnauthorized, statusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			encryptedToken := token.Encrypted{}
			err = json.Unmarshal(body, &encryptedToken)

			require.NoError(t, err)
			require.Equal(t, "", encryptedToken.Token)

			err = resp.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestVerifyHandlerSuccess(t *testing.T) {
	log := logger.NewEmpty()

	db, err := inmemory.NewMapDB()
	require.NoError(t, err)

	tt := testTime{time: time.Now()}
	us, err := service.NewUserService(db, &tt)
	require.NoError(t, err)

	h := handlers.New(log, us, us)

	creator := token.NewCreator(&tt, service.TokenExpiredTime)
	testToken, err := creator.NewJWT("login")
	require.NoError(t, err)

	etok, err := testToken.Encrypt()
	require.NoError(t, err)

	cases := []struct {
		name           string
		encryptedToken string
		token          *token.JWT
	}{
		{
			name:           "valid token 1",
			encryptedToken: etok.Token,
			token:          testToken,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest(http.MethodPost, "https://localhost:8000/verify", nil)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Authorization", "Bearer "+tc.encryptedToken)
			require.NoError(t, err)

			verifyHandler := h.Verify()
			w := httptest.NewRecorder()
			verifyHandler.ServeHTTP(w, req)

			resp := w.Result()

			sc := resp.StatusCode
			require.Equal(t, sc, http.StatusOK)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			encryptedToken := token.Encrypted{}
			err = json.Unmarshal(body, &encryptedToken)

			nt, err := creator.ParseJWT(encryptedToken.Token)
			require.NoError(t, err)
			require.Equal(t, nt.Login, tc.token.Login)
			require.True(t, nt.IssuedAt.Equal(tc.token.IssuedAt))
			require.True(t, nt.ExpiresAt.Equal(tc.token.ExpiresAt))

			err = resp.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestVerifyHandlerWrong(t *testing.T) {
	log := logger.NewEmpty()

	db, err := inmemory.NewMapDB()
	require.NoError(t, err)

	tt := testTime{time: time.Now()}
	us, err := service.NewUserService(db, &tt)
	require.NoError(t, err)

	h := handlers.New(log, us, us)

	cases := []struct {
		name           string
		encryptedToken string
		token          *token.JWT
	}{
		{
			name: "invalid token 1",
			encryptedToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
				"eyJMb2dpbiI6ImxvZ2luIiwiSXNzdWVkQXQiOiIyMDI0LTAyLTExVDE4OjU5OjIyLjExMDA0NSswMzowMCIsIkV4cGlyZXNBdCI6IjIwMjQtMDItMTFUMTk6NTk6MjIuMTEwMDQ1KzAzOjAwIn0." +
				"MBMazkvvELdU0ZvBXNeIvNPvc4BZ7GtkRSAK-0-cXSg",
			token: &token.JWT{},
		},
		{
			name:           "invalid token 2",
			encryptedToken: "",
			token:          &token.JWT{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest(http.MethodPost, "https://localhost:8000/verify", nil)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Authorization", "Bearer "+tc.encryptedToken)
			require.NoError(t, err)

			verifyHandler := h.Verify()
			w := httptest.NewRecorder()
			verifyHandler.ServeHTTP(w, req)

			resp := w.Result()

			sc := resp.StatusCode
			require.Equal(t, sc, 498)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			encryptedToken := token.Encrypted{}
			err = json.Unmarshal(body, &encryptedToken)

			require.NoError(t, err)
			require.Equal(t, "", encryptedToken.Token)

			err = resp.Body.Close()
			require.NoError(t, err)
		})
	}
}
