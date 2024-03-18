package delete_test

import (
	"errors"
	"fmt"
	"github.com/rmntim/movielab/internal/lib/logger/handlers/slogdiscard"
	actorsDelete "github.com/rmntim/movielab/internal/server/handlers/actors/delete"
	"github.com/rmntim/movielab/internal/server/handlers/actors/delete/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestActorsDelete(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		role      string
		respCode  int
		respError string
		mockError error
	}{
		{
			name:     "Success",
			id:       "1",
			respCode: http.StatusOK,
		},
		{
			name:      "Bad id",
			id:        "a",
			respCode:  http.StatusBadRequest,
			respError: "Failed to parse actor id",
		},
		{
			name:      "Unauthorized",
			id:        "1",
			role:      "user",
			respCode:  http.StatusUnauthorized,
			respError: "Insufficient permissions",
		},
		{
			name:      "DeleteActor Error",
			id:        "1",
			respCode:  http.StatusInternalServerError,
			respError: "Failed to delete actor",
			mockError: errors.New("failed to delete actor"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actorsDeleterMock := mocks.NewActorDeleter(t)

			if tt.respError == "" || tt.mockError != nil {
				actorsDeleterMock.
					On("DeleteActor", mock.AnythingOfType("int")).
					Return(tt.mockError).
					Once()
			}

			handler := actorsDelete.New(slogdiscard.NewDiscardLogger(), actorsDeleterMock)

			mux := http.NewServeMux()
			mux.HandleFunc("DELETE /{id}", handler)

			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/%s", tt.id), nil)
			require.NoError(t, err)

			role := "admin"
			if tt.role != "" {
				role = tt.role
			}
			req.Header.Set("x-role", role)

			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			require.Equal(t, tt.respCode, rr.Code)
		})
	}
}
