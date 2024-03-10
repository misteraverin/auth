package handlers

import (
	"auth/internal/domain/token"
	"auth/internal/domain/user"
	"auth/internal/errdomain"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

const TokenInvalid int = 498

type Handlers struct {
	log           *slog.Logger
	userUpdater   user.Updater
	userValidator user.Validator
}

type errorStruct struct {
	Error string `json:error`
}

func New(log *slog.Logger, userUpdater user.Updater, userValidator user.Validator) *Handlers {
	return &Handlers{log: log, userUpdater: userUpdater, userValidator: userValidator}
}

func (h *Handlers) writeTokenResponse(log *slog.Logger, w http.ResponseWriter, err error, encrypted *token.Encrypted) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case errors.Is(err, errdomain.ErrNotCorrectBasicAuth):
		status := http.StatusUnauthorized
		msg := errorStruct{Error: errdomain.ErrNotCorrectBasicAuth.Error()}
		h.writeJson(log, w, status, msg)
	case errors.Is(err, errdomain.ErrWrongLoginOrPassword):
		status := http.StatusUnauthorized
		msg := errorStruct{Error: errdomain.ErrWrongLoginOrPassword.Error()}
		h.writeJson(log, w, status, msg)
	case errors.Is(err, errdomain.ErrUserIsNotExist):
		status := http.StatusNotFound
		msg := errorStruct{Error: errdomain.ErrUserIsNotExist.Error()}
		h.writeJson(log, w, status, msg)
	case errors.Is(err, errdomain.ErrTokenExpired):
		status := TokenInvalid
		msg := errorStruct{Error: errdomain.ErrTokenExpired.Error()}
		h.writeJson(log, w, status, msg)
	case errors.Is(err, errdomain.ErrTokenInvalid):
		status := TokenInvalid
		msg := errorStruct{Error: errdomain.ErrTokenInvalid.Error()}
		h.writeJson(log, w, status, msg)
	case err != nil:
		log.Error("unexpected error: ", "error: ", err)
		status := http.StatusInternalServerError
		h.writeJson(log, w, status, "")
	default:
		status := http.StatusOK
		h.writeJson(log, w, status, encrypted)
	}
}

func (h *Handlers) writeJson(log *slog.Logger, w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Error("write response", err, status, v)
	}
}
