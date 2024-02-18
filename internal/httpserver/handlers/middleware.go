package handlers

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

const requestId = "requestId"

func AddRequestId(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uniqueId := uuid.New().String()
		log.Info("add request id", requestId, uniqueId)

		ctx := context.WithValue(r.Context(), requestId, uniqueId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getRequestId(ctx context.Context) string {
	requestId := ctx.Value(requestId)

	switch requestId.(type) {
	case string:
		return requestId.(string)
	}

	return ""
}
