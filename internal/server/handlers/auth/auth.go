package auth

import (
	"errors"
	"github.com/go-chi/render"
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

type UserRoleGetter interface {
	GetUserRole(username string, password string) (string, error)
}

func New(log *slog.Logger, userGetter UserRoleGetter, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.New"

		log := log.With(slog.String("op", op))

		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("Failed to decode request", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("Invalid request"))
			return
		}

		log.Info("Request decoded")

		if err := validator.New().Struct(req); err != nil {
			log.Error("Invalid request", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			var validationErr validator.ValidationErrors
			errors.As(err, &validationErr)
			render.JSON(w, r, resp.ValidationError(validationErr))
			return
		}

		role, err := userGetter.GetUserRole(req.Username, req.Password)
		if err != nil {
			log.Error("No user found", sl.Err(err))
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Error("No user found"))
			return
		}

		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"role": role,
		})
		token, err := jwtToken.SignedString([]byte(secret))
		if err != nil {
			log.Error("Failed to sign JWT token", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("Failed to sign JWT token"))
			return
		}

		render.JSON(w, r, Response{
			Response: resp.Ok(),
			Token:    token,
		})
	}
}
