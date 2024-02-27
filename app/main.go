package main

import (
	"auth/internal/domain/user/service"
	"auth/internal/httpserver"
	"auth/internal/httpserver/handlers"
	"auth/internal/repository/inmemory"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type baseTime struct{}

func (b *baseTime) Now() time.Time {
	return time.Now()
}

func main() {
	//todo move log settings to config
	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)

	log.Info("create repository")
	db, err := inmemory.NewMapDB()
	if err != nil {
		log.Error("create repository", err)
	}

	bt := baseTime{}
	us, err := service.NewUserService(db, &bt)
	if err != nil {
		log.Error("create user service", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("create server")
	h := handlers.New(log, us, us)
	srv := httpserver.Create(h)

	go func() {
		log.Info("start server")
		err := srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Error("stop server", "err: ", err)
		}
	}()

	<-done

	//todo move magic time number to config
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info("start shutdown server")
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown server", err)
	}
}
