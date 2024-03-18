package delete

import (
	"github.com/rmntim/movielab/internal/lib/logger/sl"
	"log/slog"
	"net/http"
	"strconv"
)

type MovieDeleter interface {
	DeleteMovie(id int) error
}

func New(log *slog.Logger, movieDeleter MovieDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.movies.delete.New"

		log := log.With(slog.String("op", op))

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Error("Failed to parse movie id", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = movieDeleter.DeleteMovie(id)
		if err != nil {
			log.Error("Failed to delete movie", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
