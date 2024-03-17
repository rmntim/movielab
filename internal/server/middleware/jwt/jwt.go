package jwt

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt"
	resp "github.com/rmntim/movielab/internal/lib/api/response"
	"net/http"
)

// New creates new middleware, sets `x-role` header if user is authorized.
func New(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token != "" {
				authorized, err := isAuthorized(token, secret)
				if authorized {
					role, err := getUserRole(token, secret)
					if err != nil {
						w.WriteHeader(http.StatusUnauthorized)
						render.JSON(w, r, resp.Error(err.Error()))
						return
					}
					w.Header().Set("x-role", role)
					next.ServeHTTP(w, r)
					return
				}
				w.WriteHeader(http.StatusUnauthorized)
				render.JSON(w, r, resp.Error(err.Error()))
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, resp.Error("Unauthorized"))
		}
		return http.HandlerFunc(fn)
	}
}

func isAuthorized(reqToken string, secret string) (bool, error) {
	_, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func getUserRole(reqToken string, secret string) (string, error) {
	token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	return claims["role"].(string), nil
}
