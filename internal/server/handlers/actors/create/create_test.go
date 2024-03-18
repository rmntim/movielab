package create_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/rmntim/movielab/internal/entity"
	"github.com/rmntim/movielab/internal/lib/logger/handlers/slogdiscard"
	"github.com/rmntim/movielab/internal/server/handlers/actors/create"
	"github.com/rmntim/movielab/internal/server/handlers/actors/create/mocks"
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
		reqActor  *entity.NewActor
		role      string
		respCode  int
		respError string
		mockError error
	}{
		{
			name: "Success",
			reqActor: &entity.NewActor{
				Name:      "Test",
				Sex:       "Test",
				BirthDate: time.Now(),
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
			name: "CreateActor Error",
			reqActor: &entity.NewActor{
				Name:      "Test",
				Sex:       "Test",
				BirthDate: time.Now(),
			},
			respCode:  http.StatusInternalServerError,
			respError: "Failed to create actor",
			mockError: errors.New("failed to create actor"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actorsCreatorMock := mocks.NewActorCreator(t)

			if tt.respError == "" || tt.mockError != nil {
				actorsCreatorMock.On("CreateActor", mock.AnythingOfType("*entity.NewActor")).Return(1, tt.mockError).Once()
			}

			handler := create.New(slogdiscard.NewDiscardLogger(), actorsCreatorMock)

			input, err := json.Marshal(tt.reqActor)
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
