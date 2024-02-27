package user

import (
	"auth/internal/domain/token"
	"context"
	"log/slog"
)

type User struct {
	Login    string
	Password string
}

type Updater interface {
	UpdateUser(ctx context.Context, log *slog.Logger, login string, password string) (*token.Encrypted, error)
}

type Validator interface {
	VerifyUser(ctx context.Context, log *slog.Logger, tokenStr string) (*token.Encrypted, error)
}
