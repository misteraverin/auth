package handlers

import (
	"BasicAuth/internal/httpserver/token"
	"BasicAuth/internal/statuscode"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

func Verify(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Verify"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", getRequestId(r.Context())),
		)

		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")

		if len(splitToken) < 2 {
			writeTokenInvalid(log, w)
			return
		}

		tokenStr := splitToken[1]

		log.Info("parse token")
		myToken, err := token.ParseJWT(tokenStr)
		if err != nil {
			slog.Error("parse token", "", err)
			writeTokenInvalid(log, w)
			return
		}

		log.Info("update token")
		encryptedToken, err := myToken.Update(TokenExpiredTime)
		if err != nil {
			writeJson(log, w, http.StatusInternalServerError, "")
			log.Error("update token", err)
			return
		}

		writeOk(log, w, encryptedToken)
	}
}

func writeTokenInvalid(log *slog.Logger, w http.ResponseWriter) {
	writeJson(
		log,
		w,
		statuscode.TokenInvalid,
		fmt.Sprintf("%s", "error: expired or otherwise invalid token"),
	)
}

func writeOk(log *slog.Logger, w http.ResponseWriter, v any) {
	writeJson(
		log,
		w,
		http.StatusOK,
		v,
	)
}
