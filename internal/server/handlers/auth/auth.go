package auth

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	resp "github.com/rmntim/movielab/internal/lib/api/response"
	"github.com/rmntim/movielab/internal/lib/logger/sl"
	"log/slog"
	"net/http"
)

type Request struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	resp.Response
	Token string `json:"token"`
}

type UserGetter interface {
	GetUserRole(username string, password string) (string, error)
}

func New(log *slog.Logger, userGetter UserGetter, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.New"

		log = log.With(slog.String("op", op))

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("Failed to decode request", sl.Err(err))
			if err := json.NewEncoder(w).Encode(resp.Error("Failed to decode request")); err != nil {
				log.Error("Failed to encode response", sl.Err(err))
			}
			return
		}

		log.Info("Request decoded")

		if err := validator.New().Struct(req); err != nil {
			log.Error("Invalid request", sl.Err(err))
			var validationErr validator.ValidationErrors
			errors.As(err, &validationErr)
			if err := json.NewEncoder(w).Encode(resp.ValidationError(validationErr)); err != nil {
				log.Error("Failed to encode response", sl.Err(err))
			}
			return
		}

		role, err := userGetter.GetUserRole(req.Username, req.Password)
		if err != nil {
			log.Error("No user found", sl.Err(err))
			if err := json.NewEncoder(w).Encode(resp.Error("No user found")); err != nil {
				log.Error("Failed to encode response", sl.Err(err))
			}
			return
		}

		log.Info("User signed in")
		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"role": role,
		})
		token, err := jwtToken.SignedString([]byte(secret))
		if err != nil {
			log.Error("Failed to sign JWT token", sl.Err(err))
			if err := json.NewEncoder(w).Encode(resp.Error("Failed to sign JWT token")); err != nil {
				log.Error("Failed to encode response", sl.Err(err))
			}
			return
		}

		if err := json.NewEncoder(w).Encode(Response{
			Response: resp.Ok(),
			Token:    token,
		}); err != nil {
			log.Error("Failed to return JWT token", sl.Err(err))
		}
	}
}
