package delete

import (
	"github.com/go-chi/render"
	resp "github.com/rmntim/movielab/internal/lib/api/response"
	"github.com/rmntim/movielab/internal/lib/logger/sl"
	"log/slog"
	"net/http"
	"strconv"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=MovieDeleter
type MovieDeleter interface {
	DeleteMovie(id int) error
}

func New(log *slog.Logger, movieDeleter MovieDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.movies.delete.New"

		log := log.With(slog.String("op", op))

		if r.Header.Get("x-role") != "admin" {
			log.Error("Insufficient permissions", slog.String("role", r.Header.Get("x-role")))
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, resp.Error("Insufficient permissions"))
			return
		}

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Error("Failed to parse movie id", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("Failed to parse movie id"))
			return
		}

		err = movieDeleter.DeleteMovie(id)
		if err != nil {
			log.Error("Failed to delete movie", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("Failed to delete movie"))
			return
		}

		render.JSON(w, r, resp.Ok())
	}
}
