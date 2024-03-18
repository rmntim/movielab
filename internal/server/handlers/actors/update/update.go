package update

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/rmntim/movielab/internal/entity"
	resp "github.com/rmntim/movielab/internal/lib/api/response"
	"github.com/rmntim/movielab/internal/lib/logger/sl"
	"github.com/rmntim/movielab/internal/storage"
	"log/slog"
	"net/http"
	"strconv"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=ActorUpdater
type ActorUpdater interface {
	GetActorById(id int) (*entity.Actor, error)
	UpdateActor(id int, actor *entity.Actor) error
}

type Response struct {
	*entity.Actor
	resp.Response
}

func New(log *slog.Logger, actorUpdater ActorUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.actors.update.New"

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

		oldActor, err := actorUpdater.GetActorById(id)
		if err != nil {
			if errors.Is(err, storage.ErrActorNotFound) {
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("Actor not found"))
				return
			}
			log.Error("Failed to get actor", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("Failed to get actor"))
			return
		}

		newActor := *oldActor
		if err := render.DecodeJSON(r.Body, &newActor); err != nil {
			log.Error("Failed to parse body", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("Failed to parse body"))
			return
		}

		if err := actorUpdater.UpdateActor(id, &newActor); err != nil {
			log.Error("Failed to update actor", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("Failed to update actor"))
			return
		}

		render.JSON(w, r, Response{
			Response: resp.Ok(),
			Actor:    &newActor,
		})
	}
}
