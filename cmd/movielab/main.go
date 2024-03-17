package main

import (
	"github.com/rmntim/movielab/internal/config"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	// init config
	cfg := config.MustLoad()

	// init logger
	logger := setupLogger(cfg.Env)

	logger.Info("Starting server", slog.String("env", cfg.Env))
	logger.Debug("Debug messages are enabled")

	// init storage

	// init router

	// start server
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
