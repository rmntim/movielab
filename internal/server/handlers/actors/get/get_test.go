package get_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rmntim/movielab/internal/entity"
	"github.com/rmntim/movielab/internal/lib/logger/handlers/slogdiscard"
	"github.com/rmntim/movielab/internal/server/handlers/actors/get"
	"github.com/rmntim/movielab/internal/server/handlers/actors/get/mocks"
	"github.com/rmntim/movielab/internal/storage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestActorsGet(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		respBody  *entity.Actor
		respCode  int
		respError string
		mockError error
	}{
		{
			name:     "Success",
			id:       "1",
			respBody: &entity.Actor{},
			respCode: http.StatusOK,
		},
		{
			name:      "Bad id",
			id:        "a",
			respCode:  http.StatusBadRequest,
			respError: "Failed to parse id",
		},
		{
			name:      "GetActorById error",
			id:        "1",
			respCode:  http.StatusInternalServerError,
			respError: "Failed to get actor",
			mockError: errors.New("failed to get actor"),
		},
		{
			name:      "Actor not found error",
			id:        "1",
			respCode:  http.StatusNotFound,
			respError: "Actor not found",
			mockError: storage.ErrActorNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actorsByIdGetterMock := mocks.NewActorByIdGetter(t)

			if tt.respError == "" || tt.mockError != nil {
				actorsByIdGetterMock.
					On("GetActorById", mock.AnythingOfType("int")).
					Return(tt.respBody, tt.mockError).
					Once()
			}

			handler := get.New(slogdiscard.NewDiscardLogger(), actorsByIdGetterMock)

			mux := http.NewServeMux()
			mux.HandleFunc("/{id}", handler)

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tt.id), nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			require.Equal(t, tt.respCode, rr.Code)
			var resp get.Response
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))

			require.Equal(t, tt.respBody, resp.Actor)
			require.Equal(t, tt.respError, resp.Error)
		})
	}
}
