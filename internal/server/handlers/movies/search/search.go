package search

import (
	"github.com/go-chi/render"
	"github.com/rmntim/movielab/internal/entity"
	"github.com/rmntim/movielab/internal/lib/api/response"
	"github.com/rmntim/movielab/internal/lib/logger/sl"
	"log/slog"
	"net/http"
	"strconv"
)

type MovieSearcher interface {
	SearchMovies(title, actorName string, limit, offset int) ([]entity.Movie, error)
}

type Response struct {
	response.Response
	Movies []entity.Movie `json:"movies"`
}

func New(log *slog.Logger, movieSearcher MovieSearcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.movies.query.New"

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
				render.JSON(w, r, response.Error("Failed to parse limit"))
				return
			}
		}
		queryOffset := r.URL.Query().Get("offset")
		if queryOffset != "" {
			offset, err = strconv.Atoi(queryOffset)
			if err != nil {
				log.Error("Failed to parse offset", sl.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Failed to parse offset"))
				return
			}
		}

		title := r.URL.Query().Get("title")
		actorName := r.URL.Query().Get("actor")

		movies, err := movieSearcher.SearchMovies(title, actorName, limit, offset)
		if err != nil {
			log.Error("Failed to search movies", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("Failed to search movies"))
			return
		}

		if movies == nil {
			movies = []entity.Movie{}
		}

		render.JSON(w, r, Response{
			Response: response.Ok(),
			Movies:   movies,
		})
	}
}
