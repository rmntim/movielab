package jwt_test

import (
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt"
	resp "github.com/rmntim/movielab/internal/lib/api/response"
	jwtMw "github.com/rmntim/movielab/internal/server/middleware/jwt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJwtNew(t *testing.T) {
	const jwtSecret = "secret"

	tests := []struct {
		name       string
		token      string
		respStatus int
		respError  string
	}{
		{
			name:       "Success",
			token:      generateJwt(t, "admin", "admin", jwtSecret),
			respStatus: http.StatusOK,
		},
		{
			name:       "Error",
			respStatus: http.StatusUnauthorized,
			respError:  "Unauthorized",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			middleware := jwtMw.New(jwtSecret)

			req, err := http.NewRequest("GET", "/", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", tt.token)

			rr := httptest.NewRecorder()

			handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				render.JSON(w, r, resp.Ok())
			}))

			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.respStatus, rr.Code)
			var res resp.Response
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &res))
			require.Equal(t, tt.respError, res.Error)
		})
	}
}

func generateJwt(t *testing.T, username, role, secret string) string {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"role":     role,
	})
	token, err := jwtToken.SignedString([]byte(secret))
	require.NoError(t, err)
	return token
}
