package handlers

import (
	"auth/internal/errdomain"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

func (h *Handlers) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Login"

		log := h.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		login, password, ok := r.BasicAuth()

		if !ok {
			h.writeTokenResponse(log, w, errdomain.ErrNotCorrectBasicAuth, nil)
			return
		}

		encryptedToken, err := h.userUpdater.UpdateUser(r.Context(), log, login, password)
		h.writeTokenResponse(log, w, err, encryptedToken)
	}
}
