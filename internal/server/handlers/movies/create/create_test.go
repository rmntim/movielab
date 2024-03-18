package create_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/rmntim/movielab/internal/entity"
	"github.com/rmntim/movielab/internal/lib/logger/handlers/slogdiscard"
	"github.com/rmntim/movielab/internal/server/handlers/movies/create"
	"github.com/rmntim/movielab/internal/server/handlers/movies/create/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMovieCreate(t *testing.T) {
	tests := []struct {
		name      string
		reqMovie  *entity.NewMovie
		role      string
		respCode  int
		respError string
		mockError error
	}{
		{
			name: "Success",
			reqMovie: &entity.NewMovie{
				Title:       "Test",
				Description: "Test",
				ReleaseDate: time.Now(),
				Rating:      1,
				ActorIDs:    []int32{1},
			},
			respCode: http.StatusOK,
		},
		{
			name:      "Unauthorized",
			role:      "user",
			respCode:  http.StatusUnauthorized,
			respError: "Insufficient permissions",
		},
		{
			name: "CreateMovie Error",
			reqMovie: &entity.NewMovie{
				Title:       "Test",
				Description: "Test",
				ReleaseDate: time.Now(),
				Rating:      1,
				ActorIDs:    []int32{1},
			},
			respCode:  http.StatusInternalServerError,
			respError: "Failed to create movie",
			mockError: errors.New("failed to create movie"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			movieCreatorMock := mocks.NewMovieCreator(t)

			if tt.respError == "" || tt.mockError != nil {
				movieCreatorMock.On("CreateMovie", mock.AnythingOfType("*entity.NewMovie")).Return(1, tt.mockError).Once()
			}

			handler := create.New(slogdiscard.NewDiscardLogger(), movieCreatorMock)

			input, err := json.Marshal(tt.reqMovie)
			require.NoError(t, err)

			mux := http.NewServeMux()
			mux.HandleFunc("POST /", handler)

			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(input))
			require.NoError(t, err)

			role := "admin"
			if tt.role != "" {
				role = tt.role
			}
			req.Header.Set("x-role", role)

			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			require.Equal(t, tt.respCode, rr.Code)
			var resp create.Response
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))

			require.Equal(t, tt.respError, resp.Error)
		})
	}
}
