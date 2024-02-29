package main

import (
	"auth/internal/database/postgres"
	"auth/models"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
	"time"
)

const secretKey = "your-secret-key"

func main() {
	authService, err := postgres.NewAuthService()
	if err != nil {
		fmt.Println("Failed to initialize database:", err)
		return
	}

	http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		Verify(w, r, authService)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		Login(w, r, authService)
	})

	http.ListenAndServe(":8080", nil)
}

func Login(w http.ResponseWriter, r *http.Request, authService *postgres.AuthService) {
	username, password, ok := r.BasicAuth()
	if !ok || username == "" || password == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	if user, err := authService.GetUserByUserName(username); err == nil {
		if user.Password == password {
			token, err := GenerateToken(authService, *user)
			if err != nil {

			}
			w.Header().Set("JWT", token)
			return
		}
	}

	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func Verify(w http.ResponseWriter, r *http.Request, authService *postgres.AuthService) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		http.Error(w, "Unauthorized", http.StatusInternalServerError)
	}

	authorizationHeaderWithoutBearer := strings.TrimSpace(strings.TrimPrefix(authorizationHeader, "Bearer"))

	token, err := jwt.Parse(authorizationHeaderWithoutBearer, func(token *jwt.Token) (interface{}, error) {
		// Здесь нужно вернуть секретный ключ, который использовался для подписи токена
		return []byte(secretKey), nil
	})
	if err != nil {

	}

	if token.Valid {
		userName, err := authService.GetUserNameByToken(token.Raw)
		if err != nil {
		}
		user, err := authService.GetUserByUserName(userName)
		if err != nil {
		}
		token, err := GenerateToken(authService, *user)
		if err != nil {

		}
		w.Header().Set("JWT", token)
		fmt.Println("Токен перегенирирован")
		return
	} else {
		http.Error(w, "Unauthorized", http.StatusInternalServerError)
	}
}

func GenerateToken(authService *postgres.AuthService, user models.User) (string, error) {
	claims := jwt.MapClaims{}
	claims["id"] = user.ID
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour).Unix() // Время истечения токена
	claims["iat"] = time.Now().Unix()                // Время создания токена

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte(secretKey)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println("Ошибка подписи токена:", err)
		return "", err
	}

	err = authService.SaveToken(signedToken, user.Username)
	if err != nil {
		fmt.Println("Ошибка подписи токена:", err)
		return "", err
	}

	return signedToken, nil
}
