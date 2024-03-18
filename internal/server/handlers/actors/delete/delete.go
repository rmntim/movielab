package delete

import (
	"github.com/go-chi/render"
	resp "github.com/rmntim/movielab/internal/lib/api/response"
	"github.com/rmntim/movielab/internal/lib/logger/sl"
	"log/slog"
	"net/http"
	"strconv"
)

type ActorDeleter interface {
	DeleteActor(id int) error
}

func New(log *slog.Logger, actorDeleter ActorDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.actors.delete.New"

		log := log.With(slog.String("op", op))

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Error("Failed to parse actor id", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = actorDeleter.DeleteActor(id)
		if err != nil {
			log.Error("Failed to delete actor", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, resp.Ok())
	}
}
