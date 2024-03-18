package auth_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/rmntim/movielab/internal/lib/logger/handlers/slogdiscard"
	"github.com/rmntim/movielab/internal/server/handlers/auth"
	"github.com/rmntim/movielab/internal/server/handlers/auth/mocks"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuth(t *testing.T) {
	const testSecret = "testsecret"

	tests := []struct {
		name      string
		username  string
		password  string
		respToken string
		respCode  int
		respError string
		mockError error
	}{
		{
			name:      "Success",
			username:  "Successful",
			password:  "Successful",
			respToken: generateJwt(t, testSecret),
			respCode:  http.StatusOK,
		},
		{
			name:      "Invalid credentials",
			username:  "Unsuccessful",
			password:  "Unsuccessful",
			respCode:  http.StatusNotFound,
			respError: "No user found",
			mockError: errors.New("no user found"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authMock := mocks.NewUserRoleGetter(t)

			if tt.mockError != nil {
				authMock.
					On("GetUserRole", "Unsuccessful", tt.password).
					Return("", tt.mockError).Once()
			} else {
				authMock.
					On("GetUserRole", "Successful", tt.password).
					Return("admin", nil).Once()
			}

			handler := auth.New(slogdiscard.NewDiscardLogger(), authMock, testSecret)

			input := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, tt.username, tt.password)

			mux := http.NewServeMux()
			mux.Handle("POST /", handler)

			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			require.Equal(t, tt.respCode, rr.Code)
			var resp auth.Response
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))

			require.Equal(t, tt.respError, resp.Error)
			require.Equal(t, tt.respToken, resp.Token)
		})
	}
}

func generateJwt(t *testing.T, secret string) string {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "admin",
	})
	token, err := jwtToken.SignedString([]byte(secret))
	require.NoError(t, err)
	return token
}
