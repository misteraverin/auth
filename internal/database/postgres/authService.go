package postgres

import (
	"auth/models"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

const (
	salt            = "salt"
	secretKeyForJWT = "your-secret-key"
	tokenTTL        = time.Hour
)

type AuthService struct {
	Repository Repository
}

func NewAuthService() (*AuthService, error) {
	repository, err := NewPostgreSQLRepository()
	if err != nil {
		fmt.Println("Failed to initialize database:", err)
		return nil, err
	}
	return &AuthService{Repository: repository}, nil
}

func (s *AuthService) GetUserByUserName(username string) (*models.User, error) {
	return s.Repository.GetUserByUserName(username)
}
func (s *AuthService) GetUserNameByToken(token string) (string, error) {
	return s.Repository.GetUserNameByToken(token)
}

func (s *AuthService) SaveToken(token, username string) error {
	userHasToken, err := s.Repository.CheckUserHasToken(username)
	if err != nil {
		return err
	}
	if userHasToken {
		s.Repository.UpdateToken(token, username)
	} else {
		s.Repository.SaveToken(token, username)
	}
	return nil
}

//	func (s *AuthService) SaveUser(user *models.User) error {
//		return s.Repository.SaveUser(user)
//	}
//
//	func (s *AuthService) GetUserByID(id int) (*models.User, error) {
//		return s.Repository.GetUserByID(id)
//	}

func (s *AuthService) Login(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok || username == "" || password == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if user, err := s.GetUserByUserName(username); err == nil {
		err = bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(password+salt))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := GenerateToken(s, *user)
		if err != nil {

		}
		w.Header().Set("JWT", token)
		return
	}

	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func (s *AuthService) Verify(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	authorizationHeaderWithoutBearer := strings.TrimSpace(strings.TrimPrefix(authorizationHeader, "Bearer"))

	token, err := jwt.Parse(authorizationHeaderWithoutBearer, func(token *jwt.Token) (interface{}, error) {
		// Здесь нужно вернуть секретный ключ, который использовался для подписи токена
		return []byte(secretKeyForJWT), nil
	})
	if err != nil {

	}

	if token.Valid {
		userName, err := s.GetUserNameByToken(token.Raw)
		if err != nil {
			fmt.Println("Нет пользователя с таким токеном")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := s.GetUserByUserName(userName)
		if err != nil {
		}
		token, err := GenerateToken(s, *user)
		if err != nil {

		}
		w.Header().Set("JWT", token)
		fmt.Println("Токен перегенирирован")
		return
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func GenerateToken(authService *AuthService, user models.User) (string, error) {
	claims := jwt.MapClaims{}
	claims["id"] = user.ID
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(tokenTTL).Unix() // Время истечения токена
	claims["iat"] = time.Now().Unix()               // Время создания токена

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte(secretKeyForJWT)
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
