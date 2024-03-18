package get_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rmntim/movielab/internal/entity"
	"github.com/rmntim/movielab/internal/lib/logger/handlers/slogdiscard"
	"github.com/rmntim/movielab/internal/server/handlers/movies/get"
	"github.com/rmntim/movielab/internal/server/handlers/movies/get/mocks"
	"github.com/rmntim/movielab/internal/storage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMoviesGet(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		respBody  *entity.Movie
		respCode  int
		respError string
		mockError error
	}{
		{
			name:     "Success",
			id:       "1",
			respBody: &entity.Movie{},
			respCode: http.StatusOK,
		},
		{
			name:      "Bad id",
			id:        "a",
			respCode:  http.StatusInternalServerError,
			respError: "Failed to parse id",
		},
		{
			name:      "GetMovieById error",
			id:        "1",
			respCode:  http.StatusInternalServerError,
			respError: "Failed to get movie",
			mockError: errors.New("failed to get movie"),
		},
		{
			name:      "Movie not found error",
			id:        "1",
			respCode:  http.StatusNotFound,
			respError: "Movie not found",
			mockError: storage.ErrMovieNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			moviesByIdGetterMock := mocks.NewMovieByIdGetter(t)

			if tt.respError == "" || tt.mockError != nil {
				moviesByIdGetterMock.
					On("GetMovieById", mock.AnythingOfType("int")).
					Return(tt.respBody, tt.mockError).
					Once()
			}

			handler := get.New(slogdiscard.NewDiscardLogger(), moviesByIdGetterMock)

			mux := http.NewServeMux()
			mux.HandleFunc("/{id}", handler)

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tt.id), nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			require.Equal(t, tt.respCode, rr.Code)
			var resp get.Response
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))

			require.Equal(t, tt.respBody, resp.Movie)
			require.Equal(t, tt.respError, resp.Error)
		})
	}
}
