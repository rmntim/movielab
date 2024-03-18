package main

import (
	"github.com/hobord/routegroup"
	"github.com/rmntim/movielab/internal/config"
	"github.com/rmntim/movielab/internal/lib/logger/sl"
	"github.com/rmntim/movielab/internal/server/handlers/auth"
	"github.com/rmntim/movielab/internal/server/handlers/movies"
	jwtMw "github.com/rmntim/movielab/internal/server/middleware/jwt"
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

const jwtSecret = "secret"

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

	handler := setupHandler(log, storage)

	log.Info("Starting server", slog.String("address", cfg.Address))
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      handler,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed to start server")
	}

	log.Info("Server stopped")
}

func setupHandler(log *slog.Logger, storage *postgres.Storage) http.Handler {
	mux := http.NewServeMux()
	root := routegroup.NewGroup(routegroup.WithMux(mux))

	root.HandleFunc("POST /auth/sign-in", auth.New(log, storage, jwtSecret))

	apiGroup := root.SubGroup("/api")
	apiGroup.Use(jwtMw.New(jwtSecret))

	movieGroup := apiGroup.SubGroup("/movies")
	movieGroup.HandleFunc("GET /", movies.NewQueryHandler(log, storage))

	// Have to put logger last, cause routegroup package is foolish with it
	handler := loggerMw.New(log)(root)
	return handler
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
