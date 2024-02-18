package main

import (
	"BasicAuth/internal/httpserver"
	"BasicAuth/internal/repository"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//todo move log settings to config
	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)

	log.Info("create repository")
	db, err := repository.NewMapDB()
	if err != nil {
		log.Error("create repository", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := httpserver.CreateServer(log, db)
	go func() {
		log.Info("start server")
		err := srv.ListenAndServe()
		log.Error("stop server", "err: ", err)
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
