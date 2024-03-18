package get

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

type ActorByIdGetter interface {
	GetActorById(id int) (*entity.Actor, error)
}

type Response struct {
	resp.Response
	Actor *entity.Actor `json:"actor"`
}

func New(log *slog.Logger, actorByIdGetter ActorByIdGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.actors.get.New"

		log := log.With(slog.String("op", op))

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Error("Failed to parse id", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("Failed to parse id"))
			return
		}

		actor, err := actorByIdGetter.GetActorById(id)
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

		render.JSON(w, r, Response{
			Response: resp.Ok(),
			Actor:    actor,
		})
	}
}
