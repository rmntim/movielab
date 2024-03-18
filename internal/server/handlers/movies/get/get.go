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

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=MovieByIdGetter
type MovieByIdGetter interface {
	GetMovieById(id int) (*entity.Movie, error)
}

type Response struct {
	resp.Response
	Movie *entity.Movie `json:"movie"`
}

func New(log *slog.Logger, movieByIdGetter MovieByIdGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.movies.get.New"

		log := log.With(slog.String("op", op))

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Error("Failed to parse id", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("Failed to parse id"))
			return
		}

		movie, err := movieByIdGetter.GetMovieById(id)
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

		render.JSON(w, r, Response{
			Response: resp.Ok(),
			Movie:    movie,
		})
	}
}
