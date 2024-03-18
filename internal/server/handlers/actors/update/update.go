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
