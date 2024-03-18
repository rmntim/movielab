package delete_test

import (
	"errors"
	"fmt"
	"github.com/rmntim/movielab/internal/lib/logger/handlers/slogdiscard"
	moviesDelete "github.com/rmntim/movielab/internal/server/handlers/movies/delete"
	"github.com/rmntim/movielab/internal/server/handlers/movies/delete/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMoviesDelete(t *testing.T) {
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
			respError: "Failed to parse movie id",
		},
		{
			name:      "Unauthorized",
			id:        "1",
			role:      "user",
			respCode:  http.StatusUnauthorized,
			respError: "Insufficient permissions",
		},
		{
			name:      "DeleteMovie error",
			id:        "1",
			respCode:  http.StatusInternalServerError,
			respError: "Failed to delete movie",
			mockError: errors.New("failed to delete movie"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			moviesDeleterMock := mocks.NewMovieDeleter(t)

			if tt.respError == "" || tt.mockError != nil {
				moviesDeleterMock.
					On("DeleteMovie", mock.AnythingOfType("int")).
					Return(tt.mockError).
					Once()
			}

			handler := moviesDelete.New(slogdiscard.NewDiscardLogger(), moviesDeleterMock)

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
