package handlers

import (
	"auth/internal/errdomain"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"strings"
)

func (h *Handlers) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Verify"

		log := h.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")

		if len(splitToken) < 2 {
			h.writeTokenResponse(log, w, errdomain.ErrTokenInvalid, nil)
			return
		}

		tokenStr := splitToken[1]
		encryptedToken, err := h.userValidator.VerifyUser(r.Context(), log, tokenStr)
		h.writeTokenResponse(log, w, err, encryptedToken)
	}
}
