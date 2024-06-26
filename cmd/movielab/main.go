package main

import (
	"github.com/hobord/routegroup"
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
	mux := http.NewServeMux()
	root := routegroup.NewGroup(routegroup.WithMux(mux))

	root.HandleFunc("POST /auth/sign-in", auth.New(log, storage, cfg.JwtSecret))

	apiGroup := root.SubGroup("/api")
	apiGroup.Use(jwtMw.New(cfg.JwtSecret))

	movieGroup := apiGroup.SubGroup("/movies")
	movieGroup.HandleFunc("GET /", moviesQuery.New(log, storage))
	movieGroup.HandleFunc("POST /", moviesCreate.New(log, storage))

	movieGroup.HandleFunc("GET /{id}", moviesGet.New(log, storage))
	movieGroup.HandleFunc("DELETE /{id}", moviesDelete.New(log, storage))
	movieGroup.HandleFunc("PUT /{id}", moviesUpdate.New(log, storage))
	movieGroup.HandleFunc("PATCH /{id}", moviesUpdate.New(log, storage))

	movieGroup.HandleFunc("GET /search", search.New(log, storage))

	actorGroup := apiGroup.SubGroup("/actors")
	actorGroup.HandleFunc("GET /", actorsQuery.New(log, storage))
	actorGroup.HandleFunc("POST /", actorsCreate.New(log, storage))

	actorGroup.HandleFunc("GET /{id}", actorsGet.New(log, storage))
	actorGroup.HandleFunc("DELETE /{id}", actorsDelete.New(log, storage))
	actorGroup.HandleFunc("PUT /{id}", actorsUpdate.New(log, storage))
	actorGroup.HandleFunc("PATCH /{id}", actorsUpdate.New(log, storage))

	doc := redoc.Redoc{
		SpecFile: "./api/openapi.yaml",
		SpecPath: "/openapi.yaml",
		DocsPath: "/docs",
	}
	docHandler := doc.Handler()

	root.Handle("/docs", docHandler)
	root.Handle("/openapi.yaml", docHandler)

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
