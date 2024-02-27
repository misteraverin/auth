package repository

import (
	"auth/internal/domain/user"
	"context"
)

type Interface interface {
	Exist(ctx context.Context, login string) (bool, error)
	IsCorrectPassword(ctx context.Context, login string, password string) (bool, error)
	GetUser(ctx context.Context, login string) (*user.User, error)
	Save(ctx context.Context, login string, password string) error
}
