package update_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rmntim/movielab/internal/entity"
	"github.com/rmntim/movielab/internal/lib/logger/handlers/slogdiscard"
	"github.com/rmntim/movielab/internal/server/handlers/movies/update"
	"github.com/rmntim/movielab/internal/server/handlers/movies/update/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	errBadId       = errors.New("failed to parse id")
	errMovieUpdate = errors.New("failed to update movie")
	errMovieGet    = errors.New("failed to get movie")
)

func TestMovieUpdate(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		reqMovie  *entity.Movie
		role      string
		respCode  int
		respError string
		mockError error
	}{
		{
			name: "Success",
			id:   "1",
			reqMovie: &entity.Movie{
				ID: 1,
				NewMovie: entity.NewMovie{
					Title:       "Test",
					Description: "Test",
					ReleaseDate: time.Now(),
					Rating:      1,
					ActorIDs:    []int32{1},
				},
			},
			respCode: http.StatusOK,
		},
		{
			name:      "Bad Id",
			id:        "a",
			respCode:  http.StatusBadRequest,
			respError: "Failed to parse id",
			mockError: errBadId,
		},
		{
			name:      "Unauthorized",
			id:        "1",
			role:      "user",
			respCode:  http.StatusUnauthorized,
			respError: "Insufficient permissions",
		},
		{
			name: "GetMovieById Error",
			id:   "1",
			reqMovie: &entity.Movie{
				ID: 1,
				NewMovie: entity.NewMovie{
					Title:       "Test",
					Description: "Test",
					ReleaseDate: time.Now(),
					Rating:      1,
					ActorIDs:    []int32{1},
				},
			},
			respCode:  http.StatusInternalServerError,
			respError: "Failed to get movie",
			mockError: errMovieGet,
		},
		{
			name: "UpdateMovie Error",
			id:   "1",
			reqMovie: &entity.Movie{
				ID: 1,
				NewMovie: entity.NewMovie{
					Title:       "Test",
					Description: "Test",
					ReleaseDate: time.Now(),
					Rating:      1,
					ActorIDs:    []int32{1},
				},
			},
			respCode:  http.StatusInternalServerError,
			respError: "Failed to update movie",
			mockError: errMovieUpdate,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			movieUpdaterMock := mocks.NewMovieUpdater(t)

			if tt.respError == "" || tt.mockError != nil {
				if errors.Is(tt.mockError, errMovieGet) {
					movieUpdaterMock.On("GetMovieById", mock.AnythingOfType("int")).Return(&entity.Movie{}, tt.mockError).Maybe()
				} else {
					movieUpdaterMock.On("GetMovieById", mock.AnythingOfType("int")).Return(&entity.Movie{}, nil).Maybe()
					movieUpdaterMock.On("UpdateMovie", mock.AnythingOfType("int"), mock.AnythingOfType("*entity.Movie")).Return(tt.mockError).Maybe()
				}
			}

			handler := update.New(slogdiscard.NewDiscardLogger(), movieUpdaterMock)

			input, err := json.Marshal(tt.reqMovie)
			require.NoError(t, err)

			mux := http.NewServeMux()
			mux.HandleFunc("PUT /{id}", handler)

			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/%s", tt.id), bytes.NewBuffer(input))
			require.NoError(t, err)

			role := "admin"
			if tt.role != "" {
				role = tt.role
			}
			req.Header.Set("x-role", role)
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, tt.respCode)
			body := rr.Body.String()
			var resp update.Response
			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, tt.respError, resp.Error)
		})
	}
}
