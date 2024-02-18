package repository

import (
	"BasicAuth/internal/domain/entity"
)

type Interface interface {
	Exist(login string) bool
	IsCorrectPassword(login string, password string) bool
	GetUser(login string) (entity.User, error)
	Save(login string, password string)
}
