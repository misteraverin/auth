package httpserver

import (
	"auth/internal/httpserver/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

const addr = ":8000"

func Create(handlers *handlers.Handlers) *http.Server {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Post("/login", handlers.Login())
	r.Post("/verify", handlers.Verify())

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	return srv
}
