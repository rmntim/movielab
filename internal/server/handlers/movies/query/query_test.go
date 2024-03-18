package query_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rmntim/movielab/internal/entity"
	"github.com/rmntim/movielab/internal/lib/logger/handlers/slogdiscard"
	"github.com/rmntim/movielab/internal/server/handlers/movies/query"
	"github.com/rmntim/movielab/internal/server/handlers/movies/query/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMovieQuery(t *testing.T) {
	tests := []struct {
		name      string
		limit     string
		offset    string
		orderBy   string
		respBody  []entity.Movie
		respCode  int
		respError string
		mockError error
	}{
		{
			name:     "Success asc",
			limit:    "10",
			offset:   "0",
			orderBy:  "+title",
			respBody: []entity.Movie{},
			respCode: http.StatusOK,
		},
		{
			name:     "Success desc",
			limit:    "10",
			offset:   "0",
			orderBy:  "-title",
			respBody: []entity.Movie{},
			respCode: http.StatusOK,
		},
		{
			name:      "Bad limit",
			limit:     "a",
			respCode:  http.StatusBadRequest,
			respError: "Failed to parse limit",
		},
		{
			name:      "Bad offset",
			offset:    "a",
			respCode:  http.StatusBadRequest,
			respError: "Failed to parse offset",
		},
		{
			name:      "GetMovies error",
			limit:     "10",
			offset:    "0",
			respCode:  http.StatusInternalServerError,
			respError: "Failed to get movies",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			movieGetterMock := mocks.NewMovieGetter(t)

			if tt.respError == "" || tt.mockError != nil {
				movieGetterMock.
					On("GetMovies", mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("bool")).
					Return(tt.respBody, tt.mockError)
			}

			handler := query.New(slogdiscard.NewDiscardLogger(), movieGetterMock)

			req, err := http.NewRequest(http.MethodGet,
				fmt.Sprintf("/?limit=%s&offset=%s&sort=%s", tt.limit, tt.offset, tt.orderBy),
				nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.respCode, rr.Code)
			var resp query.Response
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))

			require.Equal(t, tt.respBody, resp.Movies)
			require.Equal(t, tt.respError, resp.Error)
		})
	}
}
