package handlers_test

import (
	"BasicAuth/internal/httpserver/handlers"
	"BasicAuth/internal/httpserver/token"
	"BasicAuth/internal/statuscode"
	"BasicAuth/pkg/emptylogger"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVerifyHandler(t *testing.T) {
	tok, err := token.New("login", handlers.TokenExpiredTime)
	assert.NoError(t, err)
	etok, err := tok.Encrypt()
	assert.NoError(t, err)

	cases := []struct {
		name           string
		isValid        bool
		encryptedToken string
		token          *token.JWTToken
	}{
		{
			name:    "1st: invalid token 1",
			isValid: false,
			encryptedToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
				"eyJMb2dpbiI6ImxvZ2luIiwiSXNzdWVkQXQiOiIyMDI0LTAyLTExVDE4OjU5OjIyLjExMDA0NSswMzowMCIsIkV4cGlyZXNBdCI6IjIwMjQtMDItMTFUMTk6NTk6MjIuMTEwMDQ1KzAzOjAwIn0." +
				"MBMazkvvELdU0ZvBXNeIvNPvc4BZ7GtkRSAK-0-cXSg",
			token: &token.JWTToken{},
		},
		{
			name:           "2nd: invalid token 2",
			isValid:        false,
			encryptedToken: "",
			token:          &token.JWTToken{},
		},
		{
			name:           "3rd: valid token 1",
			isValid:        true,
			encryptedToken: etok.Token,
			token:          tok,
		},
	}

	log := emptylogger.New()
	server := httptest.NewServer(handlers.Verify(log))

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest(http.MethodPost, server.URL+"/verify", nil)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Authorization", "Bearer "+tc.encryptedToken)
			assert.NoError(t, err)

			client := server.Client()
			resp, err := client.Do(req)
			defer resp.Body.Close()

			assert.NoError(t, err)

			sc := resp.StatusCode
			if tc.isValid {
				assert.Equal(t, sc, http.StatusOK)
			} else {
				assert.Equal(t, sc, statuscode.TokenInvalid)
			}

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			encryptedToken := token.Encrypted{}
			err = json.Unmarshal(body, &encryptedToken)

			if tc.isValid {
				nt, err := token.ParseJWT(encryptedToken.Token)
				assert.NoError(t, err)
				assert.Equal(t, nt.Login, tc.token.Login)
				assert.True(t, nt.IssuedAt.After(tc.token.IssuedAt))
				assert.True(t, nt.ExpiresAt.After(tc.token.ExpiresAt))
			} else {
				assert.Error(t, err)
				assert.Equal(t, "", encryptedToken.Token)
			}
		})
	}
}
