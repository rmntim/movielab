package query_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rmntim/movielab/internal/entity"
	"github.com/rmntim/movielab/internal/lib/logger/handlers/slogdiscard"
	"github.com/rmntim/movielab/internal/server/handlers/actors/query"
	"github.com/rmntim/movielab/internal/server/handlers/actors/query/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestActorQuery(t *testing.T) {
	tests := []struct {
		name      string
		limit     string
		offset    string
		respBody  []entity.Actor
		respCode  int
		respError string
		mockError error
	}{
		{
			name:     "Success",
			limit:    "10",
			offset:   "0",
			respBody: []entity.Actor{},
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
			name:      "GetActors Error",
			limit:     "10",
			offset:    "0",
			respCode:  http.StatusInternalServerError,
			respError: "Failed to get actors",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actorGetterMock := mocks.NewActorGetter(t)

			if tt.respError == "" || tt.mockError != nil {
				actorGetterMock.
					On("GetActors", mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(tt.respBody, tt.mockError)
			}

			handler := query.New(slogdiscard.NewDiscardLogger(), actorGetterMock)

			req, err := http.NewRequest(http.MethodGet,
				fmt.Sprintf("/?limit=%s&offset=%s", tt.limit, tt.offset),
				nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.respCode, rr.Code)
			var resp query.Response
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))

			require.Equal(t, tt.respBody, resp.Actors)
			require.Equal(t, tt.respError, resp.Error)
		})
	}
}
