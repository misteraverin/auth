package tests

import (
	"auth/internal/domain/token"
	"auth/internal/domain/user/service"
	"auth/internal/httpserver/handlers"
	"auth/internal/repository/inmemory"
	"auth/pkg/logger"
	"github.com/go-chi/chi/v5"
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

func Test(t *testing.T) {
	srv := createServer(t)
	defer srv.Close()

	t.Run("valid auth", func(t *testing.T) {
		ValidAuth(t, srv)
	})

	t.Run("expired token", func(t *testing.T) {
		ExpiredToken(t, srv)
	})

	t.Run("invalid auth", func(t *testing.T) {
		InvalidAuth(t, srv)
	})

	t.Run("invalid token", func(t *testing.T) {
		InvalidToken(t, srv)
	})
}

func createServer(t *testing.T) *httptest.Server {
	r := chi.NewRouter()

	log := logger.NewEmpty()

	db, err := inmemory.NewMapDB()
	require.NoError(t, err)

	tt := testTime{time: time.Now()}
	us, err := service.NewUserService(db, &tt)
	require.NoError(t, err)

	h := handlers.New(log, us, us)

	r.Post("/login", h.Login())
	r.Post("/verify", h.Verify())

	srv := httptest.NewServer(r)

	return srv
}

func ValidAuth(t *testing.T, srv *httptest.Server) {
	testCases := []struct {
		name   string
		path   string
		method string
	}{
		{
			name:   "Valid Login 1",
			path:   "/login",
			method: http.MethodPost,
		},
		// todo: add more test cases
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			url := srv.URL + tc.path

			req, err := http.NewRequest(tc.method, url, nil)
			req.SetBasicAuth("login", "password")
			require.NoError(t, err)

			client := srv.Client()
			resp, err := client.Do(req)

			require.NotEqual(t, resp, nil)
			require.NoError(t, err)

			err = resp.Body.Close()
			require.NoError(t, err)
		})
	}
}

func InvalidToken(t *testing.T, srv *httptest.Server) {
	testCases := []struct {
		name  string
		token string
	}{
		{
			name:  "invalid token 1",
			token: "123.123.123",
		},
		// todo: add more test cases
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			url := srv.URL + "/verify"

			req, err := http.NewRequest(http.MethodPost, url, nil)
			req.Header.Set("Authorization", "Bearer "+tc.token)
			require.NoError(t, err)

			client := srv.Client()
			resp, err := client.Do(req)

			require.NotEqual(t, resp, nil)
			require.NoError(t, err)
			require.Equal(t, resp.StatusCode, 498)

			b, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			bodyStr := string(b)

			err = resp.Body.Close()
			require.NoError(t, err)

			require.Equal(t, bodyStr, "{\"Error\":\"token invalid\"}\n")
		})
	}
}

func ExpiredToken(t *testing.T, srv *httptest.Server) {
	testCases := []struct {
		name  string
		login string
	}{
		{
			name:  "Valid Login 1",
			login: "test login 1",
		},
		// todo: add more test cases
	}

	tt := &testTime{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			//add login
			url := srv.URL + "/login"

			req, err := http.NewRequest(http.MethodPost, url, nil)
			req.SetBasicAuth(tc.login, "password")
			require.NoError(t, err)

			client := srv.Client()
			resp, err := client.Do(req)

			require.NotEqual(t, resp, nil)
			require.NoError(t, err)

			err = resp.Body.Close()
			require.NoError(t, err)

			//add expired token
			tokenCreator := token.NewCreator(tt, time.Duration(0*time.Second))
			time.Sleep(1 * time.Second)

			url = srv.URL + "/verify"

			req, err = http.NewRequest(http.MethodPost, url, nil)
			require.NoError(t, err)
			token, err := tokenCreator.NewJWT(tc.login)
			require.NoError(t, err)
			encrypted, err := token.Encrypt()
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+encrypted.Token)

			resp, err = client.Do(req)

			require.NotEqual(t, resp, nil)
			require.NoError(t, err)
			require.Equal(t, resp.StatusCode, 498)

			b, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			bodyStr := string(b)

			err = resp.Body.Close()
			require.NoError(t, err)

			require.Equal(t, bodyStr, "{\"Error\":\"token expired\"}\n")
		})
	}
}

func InvalidAuth(t *testing.T, srv *httptest.Server) {
	testCases := []struct {
		name   string
		path   string
		method string
		status int
		error  string
	}{
		{
			name:   "Valid Login 1",
			path:   "/login",
			method: http.MethodGet,
			status: http.StatusMethodNotAllowed,
			error:  "",
		},
		// todo: add more test cases
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			url := srv.URL + tc.path

			req, err := http.NewRequest(tc.method, url, nil)
			req.SetBasicAuth("login", "password")
			require.NoError(t, err)

			client := srv.Client()
			resp, err := client.Do(req)
			require.NoError(t, err)
			require.Equal(t, resp.StatusCode, tc.status)
			require.NotEqual(t, resp, nil)

			err = resp.Body.Close()
			require.NoError(t, err)
		})
	}
}
