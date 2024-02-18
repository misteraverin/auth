package httpserver

import (
	"BasicAuth/internal/httpserver/handlers"
	"BasicAuth/internal/repository"
	"log/slog"
	"net/http"
)

const addr = ":8000"

func CreateServer(log *slog.Logger, rep repository.Interface) *http.Server {
	log.Info("create server")

	mux := http.NewServeMux()

	mux.HandleFunc("/login", handlers.Login(log, rep))
	mux.HandleFunc("/verify", handlers.Verify(log))

	//add middleware
	handler := handlers.AddRequestId(log, mux)

	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return srv
}
