package login

import (
	"auth/internal/models/errorsi"
	"auth/internal/models/users"
	responses "auth/internal/utils/api/response"
	"auth/internal/utils/jwt"
	"auth/internal/utils/logger/sl"

	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type Request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Responce struct {
	responses.Responce
	ID    int64  `json:"id"`
	Email string `json:"email"`
	JWT   string `json:"jwt"`
}

type userStorage interface {
	GetUserByEmail(ctx context.Context, email string) (*users.Users, error)
	IsUserNotExistsByEmail(email string) error
}

func New(log *slog.Logger, login userStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.login.new"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			http.Error(w, "failed to decode request body", http.StatusBadRequest)
			return
		}

		log.Info("request body decoded successfully", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(validateErr))
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}
		if err := login.IsUserNotExistsByEmail(req.Email); err != nil {
			msg, status := validateLoginData(err, *log)
			http.Error(w, msg, status)
			return
		}

		user, err := login.GetUserByEmail(r.Context(), req.Email)
		if err != nil {
			log.Error("failed to get user by email", sl.Err(err))
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			msg, status := validateLoginData(err, *log)
			http.Error(w, msg, status)
			return
		}

		jwtToken, err := jwt.New(int(user.ID), user.Email)
		if err != nil {
			log.Error("failed to create JWT token", sl.Err(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		log.Info("user login", slog.Int64("ID", int64(user.ID)), slog.String("email", user.Email))

		render.JSON(w, r, Responce{
			Responce: responses.Ok(),
			ID:       int64(user.ID),
			Email:    user.Email,
			JWT:      jwtToken,
		})
	}
}

func validateLoginData(err error, log slog.Logger) (string, int) {
	switch err {
	case errorsi.ErrEmailNotExists:
		log.Error("Email not exists", sl.Err(err))
		return "Email or password is incorrect", http.StatusUnauthorized
	case bcrypt.ErrMismatchedHashAndPassword:
		log.Error("Invalid password", sl.Err(err))
		return "Email or password is incorrect", http.StatusUnauthorized
	default:
		log.Error("Internal server error", sl.Err(err))
		return "Internal server error", http.StatusInternalServerError
	}
}
