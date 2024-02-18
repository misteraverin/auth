package handlers

import (
	"BasicAuth/internal/errdomain"
	"BasicAuth/internal/httpserver/token"
	"BasicAuth/internal/repository"
	"fmt"
	"log/slog"
	"net/http"
)

func Login(log *slog.Logger, rep repository.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Login"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", getRequestId(r.Context())),
		)

		login, password, ok := r.BasicAuth()

		if ok {
			if !rep.Exist(login) {
				rep.Save(login, password)
			} else if !rep.IsCorrectPassword(login, password) {
				writeWrongLoginOrPassword(log, w)
				return
			}
		} else {
			writeWrongLoginOrPassword(log, w)
			return
		}

		log.Info("create or update token")
		newToken, err := token.New(login, TokenExpiredTime)
		if err != nil {
			log.Error("create or update token", err)
		}

		log.Info("encrypt token")
		encryptedToken, err := newToken.Encrypt()
		if err != nil {
			writeJson(log, w, http.StatusInternalServerError, "")
			log.Error("encrypt token", err)
			return
		}

		writeOk(log, w, encryptedToken)
	}
}

func writeWrongLoginOrPassword(log *slog.Logger, w http.ResponseWriter) {
	writeJson(
		log,
		w,
		http.StatusUnauthorized,
		fmt.Sprintf("{ %s %s }", "error: ", errdomain.ErrWrongLoginOrPassword),
	)
}
