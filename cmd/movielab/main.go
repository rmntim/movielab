package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/mvrilo/go-redoc"
	"github.com/rmntim/movielab/internal/config"
	"github.com/rmntim/movielab/internal/lib/logger/sl"
	actorsCreate "github.com/rmntim/movielab/internal/server/handlers/actors/create"
	actorsDelete "github.com/rmntim/movielab/internal/server/handlers/actors/delete"
	actorsGet "github.com/rmntim/movielab/internal/server/handlers/actors/get"
	actorsQuery "github.com/rmntim/movielab/internal/server/handlers/actors/query"
	actorsUpdate "github.com/rmntim/movielab/internal/server/handlers/actors/update"
	"github.com/rmntim/movielab/internal/server/handlers/auth"
	moviesCreate "github.com/rmntim/movielab/internal/server/handlers/movies/create"
	moviesDelete "github.com/rmntim/movielab/internal/server/handlers/movies/delete"
	moviesGet "github.com/rmntim/movielab/internal/server/handlers/movies/get"
	moviesQuery "github.com/rmntim/movielab/internal/server/handlers/movies/query"
	"github.com/rmntim/movielab/internal/server/handlers/movies/search"
	moviesUpdate "github.com/rmntim/movielab/internal/server/handlers/movies/update"
	jwtMw "github.com/rmntim/movielab/internal/server/middleware/jwt"
	loggerMw "github.com/rmntim/movielab/internal/server/middleware/logger"
	"github.com/rmntim/movielab/internal/storage/postgres"
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

	handler := setupHandler(cfg, log, storage)

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

func setupHandler(cfg *config.Config, log *slog.Logger, storage *postgres.Storage) http.Handler {
	router := chi.NewRouter()
	router.Use(loggerMw.New(log))

	router.Post("/auth/sign-in", auth.New(log, storage, cfg.JwtSecret))

	router.Route("/api/v1", func(r chi.Router) {
		r.Use(jwtMw.New(cfg.JwtSecret))

		r.Route("/movies", setupMovieHandler(log, storage))
		r.Route("/actors", setupActorHandler(log, storage))
	})

	doc := redoc.Redoc{
		SpecFile: "./api/openapi.yaml",
		SpecPath: "/openapi.yaml",
		DocsPath: "/docs",
	}
	docHandler := doc.Handler()

	router.Handle("/docs", docHandler)
	router.Handle("/openapi.yaml", docHandler)

	return router
}

func setupMovieHandler(log *slog.Logger, storage *postgres.Storage) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/", moviesQuery.New(log, storage))
		r.Post("/", moviesCreate.New(log, storage))

		r.Get("/{id}", moviesGet.New(log, storage))
		r.Delete("/{id}", moviesDelete.New(log, storage))
		r.Put("/{id}", moviesUpdate.New(log, storage))
		r.Patch("/{id}", moviesUpdate.New(log, storage))

		r.Get("/search", search.New(log, storage))
	}
}

func setupActorHandler(log *slog.Logger, storage *postgres.Storage) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/", actorsQuery.New(log, storage))
		r.Post("/", actorsCreate.New(log, storage))

		r.Get("/{id}", actorsGet.New(log, storage))
		r.Delete("/{id}", actorsDelete.New(log, storage))
		r.Put("/{id}", actorsUpdate.New(log, storage))
		r.Patch("/{id}", actorsUpdate.New(log, storage))
	}
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
