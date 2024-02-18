package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

// todo move to config
const TokenExpiredTime = time.Duration(3600 * time.Second)

func writeJson(log *slog.Logger, w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)

	if err != nil {
		log.Error("write response", err, status, v)
	}
}
