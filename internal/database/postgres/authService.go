package postgres

import (
	"auth/models"
	"fmt"
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

func (s *AuthService) SaveUser(user *models.User) error {
	return s.Repository.SaveUser(user)
}
func (s *AuthService) GetUserByID(id int) (*models.User, error) {
	return s.Repository.GetUserByID(id)
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
