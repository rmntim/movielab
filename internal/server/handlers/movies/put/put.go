package put

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

type MovieUpdater interface {
	GetMovieById(id int) (*entity.Movie, error)
	UpdateMovie(id int, movie *entity.Movie) error
}

type Response struct {
	*entity.Movie
	resp.Response
}

func New(log *slog.Logger, movieUpdater MovieUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.movies.NewPutHandler"

		log := log.With(slog.String("op", op))

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Error("Failed to parse id", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("Failed to parse id"))
			return
		}

		oldMovie, err := movieUpdater.GetMovieById(id)
		if err != nil {
			if errors.Is(err, storage.ErrMovieNotFound) {
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("Movie not found"))
				return
			}
			log.Error("Failed to get movie", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("Failed to get movie"))
			return
		}

		newMovie := *oldMovie
		if err := render.DecodeJSON(r.Body, &newMovie); err != nil {
			log.Error("Failed to parse body", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("Failed to parse body"))
			return
		}

		if err := movieUpdater.UpdateMovie(id, &newMovie); err != nil {
			log.Error("Failed to update movie", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("Failed to update movie"))
			return
		}

		render.JSON(w, r, Response{
			Response: resp.Ok(),
			Movie:    &newMovie,
		})
	}
}
