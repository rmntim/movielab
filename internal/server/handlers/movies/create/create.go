package create

import (
	"github.com/go-chi/render"
	"github.com/rmntim/movielab/internal/entity"
	resp "github.com/rmntim/movielab/internal/lib/api/response"
	"github.com/rmntim/movielab/internal/lib/logger/sl"
	"log/slog"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=MovieCreator
type MovieCreator interface {
	CreateMovie(movie *entity.NewMovie) (int, error)
}

type Response struct {
	resp.Response
	Movie *entity.Movie `json:"movie"`
}

func New(log *slog.Logger, movieCreator MovieCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.movies.create.New"

		log := log.With(slog.String("op", op))

		if r.Header.Get("x-role") != "admin" {
			log.Error("Insufficient permissions", slog.String("role", r.Header.Get("x-role")))
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, resp.Error("Insufficient permissions"))
			return
		}

		var movie entity.NewMovie
		if err := render.DecodeJSON(r.Body, &movie); err != nil {
			log.Error("Failed to decode request", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("Invalid request"))
			return
		}

		id, err := movieCreator.CreateMovie(&movie)
		if err != nil {
			log.Error("Failed to create movie", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("Failed to create movie"))
			return
		}

		render.JSON(w, r, Response{
			Response: resp.Ok(),
			Movie:    &entity.Movie{ID: id, NewMovie: movie},
		})
	}
}
