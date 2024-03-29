package create

import (
	"github.com/go-chi/render"
	"github.com/rmntim/movielab/internal/entity"
	resp "github.com/rmntim/movielab/internal/lib/api/response"
	"github.com/rmntim/movielab/internal/lib/logger/sl"
	"log/slog"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=ActorCreator
type ActorCreator interface {
	CreateActor(actor *entity.NewActor) (int, error)
}

type Response struct {
	resp.Response
	Actor *entity.Actor `json:"actor"`
}

func New(log *slog.Logger, actorCreator ActorCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.actors.create.New"

		log := log.With(slog.String("op", op))

		if r.Header.Get("x-role") != "admin" {
			log.Error("Insufficient permissions", slog.String("role", r.Header.Get("x-role")))
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, resp.Error("Insufficient permissions"))
			return
		}

		var actor entity.NewActor
		if err := render.DecodeJSON(r.Body, &actor); err != nil {
			log.Error("Failed to decode request", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("Invalid request"))
			return
		}

		id, err := actorCreator.CreateActor(&actor)
		if err != nil {
			log.Error("Failed to create actor", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("Failed to create actor"))
			return
		}

		render.JSON(w, r, Response{
			Response: resp.Ok(),
			Actor:    &entity.Actor{ID: id, NewActor: actor},
		})
	}
}
