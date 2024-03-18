package query

import (
	"github.com/go-chi/render"
	"github.com/rmntim/movielab/internal/entity"
	resp "github.com/rmntim/movielab/internal/lib/api/response"
	"github.com/rmntim/movielab/internal/lib/logger/sl"
	"log/slog"
	"net/http"
	"strconv"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=ActorGetter
type ActorGetter interface {
	GetActors(limit, offset int) ([]entity.Actor, error)
}

type Response struct {
	resp.Response
	Actors []entity.Actor `json:"actors"`
}

func New(log *slog.Logger, actorGetter ActorGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.actors.query.New"

		log := log.With(slog.String("op", op))

		var (
			limit  = 10
			offset = 0
		)
		var err error

		queryLimit := r.URL.Query().Get("limit")
		if queryLimit != "" {
			limit, err = strconv.Atoi(queryLimit)
			if err != nil {
				log.Error("Failed to parse limit", sl.Err(err))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, resp.Error("Failed to parse limit"))
				return
			}
		}
		queryOffset := r.URL.Query().Get("offset")
		if queryOffset != "" {
			offset, err = strconv.Atoi(queryOffset)
			if err != nil {
				log.Error("Failed to parse offset", sl.Err(err))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, resp.Error("Failed to parse offset"))
				return
			}
		}

		actors, err := actorGetter.GetActors(limit, offset)
		if err != nil {
			log.Error("Failed to get actors", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("Failed to get actors"))
			return
		}

		render.JSON(w, r, Response{
			Response: resp.Ok(),
			Actors:   actors,
		})
	}
}
