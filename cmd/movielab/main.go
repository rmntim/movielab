package main

import (
	"github.com/rmntim/movielab/internal/config"
	"github.com/rmntim/movielab/internal/lib/logger/sl"
	loggerMw "github.com/rmntim/movielab/internal/server/middleware/logger"
	"github.com/rmntim/movielab/internal/storage/postgres"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Starting server", slog.String("env", cfg.Env))
	log.Debug("Debug messages are enabled")

	storage, err := postgres.New(cfg.DBUrl)
	if err != nil {
		log.Error("Failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	_ = storage

	var router http.Handler = http.NewServeMux()
	router = loggerMw.New(log)(router)

	log.Info("Starting server", slog.String("address", cfg.Address))
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed to start server")
	}

	log.Info("Server stopped")
}
func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}