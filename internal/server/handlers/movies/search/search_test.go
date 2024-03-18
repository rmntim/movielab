package search_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rmntim/movielab/internal/entity"
	"github.com/rmntim/movielab/internal/lib/logger/handlers/slogdiscard"
	"github.com/rmntim/movielab/internal/server/handlers/movies/search"
	"github.com/rmntim/movielab/internal/server/handlers/movies/search/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMovieSearch(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		actor     string
		limit     string
		offset    string
		respBody  []entity.Movie
		respCode  int
		respError string
		mockError error
	}{
		{
			name:     "Success",
			title:    "Test",
			actor:    "Test",
			limit:    "10",
			offset:   "0",
			respBody: []entity.Movie{},
			respCode: http.StatusOK,
		},
		{
			name:      "Bad limit",
			limit:     "a",
			respCode:  http.StatusInternalServerError,
			respError: "Failed to parse limit",
		},
		{
			name:      "Bad offset",
			offset:    "a",
			respCode:  http.StatusInternalServerError,
			respError: "Failed to parse offset",
		},
		{
			name:      "SearchMovies error",
			title:     "Test",
			actor:     "Test",
			limit:     "10",
			offset:    "0",
			respCode:  http.StatusInternalServerError,
			respError: "Failed to search movies",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			movieSearcherMock := mocks.NewMovieSearcher(t)

			if tt.respError == "" || tt.mockError != nil {
				movieSearcherMock.
					On("SearchMovies", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(nil, tt.mockError).
					Once()
			}

			handler := search.New(slogdiscard.NewDiscardLogger(), movieSearcherMock)

			req, err := http.NewRequest(http.MethodGet,
				fmt.Sprintf("/?title=%s&actor=%s&limit=%s&offset=%s", tt.title, tt.actor, tt.limit, tt.offset),
				nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.respCode, rr.Code)
			var resp search.Response
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))

			require.Equal(t, tt.respBody, resp.Movies)
			require.Equal(t, tt.respError, resp.Error)
		})
	}
}
