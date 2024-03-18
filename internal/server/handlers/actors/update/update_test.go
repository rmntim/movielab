package update_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rmntim/movielab/internal/entity"
	"github.com/rmntim/movielab/internal/lib/logger/handlers/slogdiscard"
	"github.com/rmntim/movielab/internal/server/handlers/actors/update"
	"github.com/rmntim/movielab/internal/server/handlers/actors/update/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	errBadId       = errors.New("failed to parse id")
	errMovieUpdate = errors.New("failed to update actor")
	errMovieGet    = errors.New("failed to get actor")
)

func TestActorUpdate(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		reqActor  *entity.Actor
		role      string
		respCode  int
		respError string
		mockError error
	}{
		{
			name: "Success",
			id:   "1",
			reqActor: &entity.Actor{
				ID: 1,
				NewActor: entity.NewActor{
					Name:      "Test",
					Sex:       "Test",
					BirthDate: time.Now(),
				},
			},
			respCode: http.StatusOK,
		},
		{
			name:      "Bad Id",
			id:        "a",
			respCode:  http.StatusBadRequest,
			respError: "Failed to parse actor id",
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
			name: "GetActorById Error",
			id:   "1",
			reqActor: &entity.Actor{
				ID: 1,
				NewActor: entity.NewActor{
					Name:      "Test",
					Sex:       "Test",
					BirthDate: time.Now(),
				},
			},
			respCode:  http.StatusInternalServerError,
			respError: "Failed to get actor",
			mockError: errMovieGet,
		},
		{
			name: "UpdateActor Error",
			id:   "1",
			reqActor: &entity.Actor{
				ID: 1,
				NewActor: entity.NewActor{
					Name:      "Test",
					Sex:       "Test",
					BirthDate: time.Now(),
				},
			},
			respCode:  http.StatusInternalServerError,
			respError: "Failed to update actor",
			mockError: errMovieUpdate,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actorUpdaterMock := mocks.NewActorUpdater(t)

			if tt.respError == "" || tt.mockError != nil {
				if errors.Is(tt.mockError, errMovieGet) {
					actorUpdaterMock.On("GetActorById", mock.AnythingOfType("int")).Return(&entity.Actor{}, tt.mockError).Maybe()
				} else {
					actorUpdaterMock.On("GetActorById", mock.AnythingOfType("int")).Return(&entity.Actor{}, nil).Maybe()
					actorUpdaterMock.On("UpdateActor", mock.AnythingOfType("int"), mock.AnythingOfType("*entity.Actor")).Return(tt.mockError).Maybe()
				}
			}

			handler := update.New(slogdiscard.NewDiscardLogger(), actorUpdaterMock)

			input, err := json.Marshal(tt.reqActor)
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
