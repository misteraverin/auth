package service

import (
	"auth/internal/domain/clock"
	"auth/internal/domain/token"
	"auth/internal/errdomain"
	"auth/internal/repository"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

// todo move to config
const TokenExpiredTime = time.Duration(3600 * time.Second)

type UserService struct {
	rep          repository.Interface
	tokenCreator *token.Creator
}

func NewUserService(rep repository.Interface, clock clock.Interface) (*UserService, error) {
	tc := token.NewCreator(clock, TokenExpiredTime)
	return &UserService{rep: rep, tokenCreator: tc}, nil
}

func (us *UserService) UpdateUser(ctx context.Context, log *slog.Logger, login string, password string) (*token.Encrypted, error) {
	ok, err := us.rep.Exist(ctx, login)

	if err != nil {
		return nil, fmt.Errorf("exist player: %w", err)
	}

	if !ok {
		err = us.rep.Save(ctx, login, password)
		if err != nil {
			return nil, fmt.Errorf("save: %w", err)
		}
	}

	ok, err = us.rep.IsCorrectPassword(ctx, login, password)

	if err != nil {
		return nil, fmt.Errorf("correct password: %w", err)
	}

	if !ok {
		return nil, errdomain.ErrWrongLoginOrPassword
	}

	log.Info("create or update token")
	newToken, err := us.tokenCreator.NewJWT(login)
	if err != nil {
		return nil, fmt.Errorf("create or update token: %w", err)
	}

	log.Info("encrypt token")
	encryptedToken, err := newToken.Encrypt()
	if err != nil {
		return nil, fmt.Errorf("encrypt token: %w", err)
	}

	return encryptedToken, nil
}

func (us *UserService) VerifyUser(ctx context.Context, log *slog.Logger, tokenStr string) (*token.Encrypted, error) {
	log.Info("parse token")
	myToken, err := us.tokenCreator.ParseJWT(tokenStr)

	if err != nil {
		if errors.Is(err, errdomain.ErrTokenInvalid) {
			return nil, errdomain.ErrTokenInvalid
		}

		return nil, fmt.Errorf("parse token: %w", err)
	}

	log.Info("update token")
	encryptedToken, err := us.tokenCreator.UpdateJWT(myToken.Login)
	if err != nil {
		if errors.Is(err, errdomain.ErrTokenExpired) {
			return nil, errdomain.ErrTokenExpired
		}

		return nil, fmt.Errorf("update token: %w", err)
	}

	return encryptedToken, nil
}
