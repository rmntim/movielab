package delete

import (
	"github.com/go-chi/render"
	resp "github.com/rmntim/movielab/internal/lib/api/response"
	"github.com/rmntim/movielab/internal/lib/logger/sl"
	"log/slog"
	"net/http"
	"strconv"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=ActorDeleter
type ActorDeleter interface {
	DeleteActor(id int) error
}

func New(log *slog.Logger, actorDeleter ActorDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.actors.delete.New"

		log := log.With(slog.String("op", op))

		if r.Header.Get("x-role") != "admin" {
			log.Error("Insufficient permissions", slog.String("role", r.Header.Get("x-role")))
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, resp.Error("Insufficient permissions"))
			return
		}

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Error("Failed to parse actor id", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("Failed to parse actor id"))
			return
		}

		err = actorDeleter.DeleteActor(id)
		if err != nil {
			log.Error("Failed to delete actor", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("Failed to delete actor"))
			return
		}

		render.JSON(w, r, resp.Ok())
	}
}
