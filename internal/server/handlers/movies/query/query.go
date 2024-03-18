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

type MovieGetter interface {
	GetMovies(limit, offset int, orderBy string, asc bool) ([]entity.Movie, error)
}

type Response struct {
	resp.Response
	Movies []entity.Movie `json:"movies"`
}

func New(log *slog.Logger, movieGetter MovieGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.movies.NewQueryHandler"

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
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, resp.Error("Failed to parse limit"))
				return
			}
		}
		queryOffset := r.URL.Query().Get("offset")
		if queryOffset != "" {
			offset, err = strconv.Atoi(queryOffset)
			if err != nil {
				log.Error("Failed to parse offset", sl.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, resp.Error("Failed to parse offset"))
				return
			}
		}

		orderBy := "title"
		asc := false

		querySort := r.URL.Query().Get("sort")
		if querySort != "" {
			if querySort[0] == '+' {
				asc = true
			}
			orderBy = querySort[1:]
		}

		movies, err := movieGetter.GetMovies(limit, offset, orderBy, asc)
		if err != nil {
			log.Error("Failed to get movies", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("Failed to get movies"))
			return
		}
		if movies == nil {
			movies = []entity.Movie{}
		}

		render.JSON(w, r, Response{
			Response: resp.Ok(),
			Movies:   movies,
		})
	}
}
