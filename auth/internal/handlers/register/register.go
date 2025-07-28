package register

import (
	"auth/internal/models/errorsi"
	resp "auth/internal/utils/api/response"
	"auth/internal/utils/logger/sl"
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

type Request struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=8"`
}

type Response struct {
	resp.Responce
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// UserStorage defines the interface for user registration storage.
type UserStorage interface {
	RegisterUser(ctx context.Context, email, username string, passHash []byte, createAt time.Time) (*Response, error)
	IsUserExistsByEmail(email string) error
	IsUserExistByUserName(username string) error
}

// User represents the user model returned by RegisterUser.
type User struct {
	ID       int64
	Email    string
	Username string
}

func New(log *slog.Logger, register UserStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.register.New"

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

		log.Info(("request body decoded successfully"), slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(validateErr))
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}
		if err := register.IsUserExistsByEmail(req.Email); err != nil {
			msg, status := validateEmailorUsername(err, *log)
			http.Error(w, msg, status)
			return
		}

		if err := register.IsUserExistByUserName(req.Username); err != nil {
			msg, status := validateEmailorUsername(err, *log)
			http.Error(w, msg, status)
			return
		}

		passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to general PassHash", sl.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		createAt := time.Now()
		res, err := register.RegisterUser(r.Context(), req.Email, req.Username, passHash, createAt)
		if err != nil {
			out, status := validateEmailorUsername(err, *log)
			http.Error(w, out, status)
			return
		}

		log.Info("user register", slog.Int64("ID", res.ID), slog.String("email", res.Email))

		render.JSON(w, r, Response{
			Responce: resp.Ok(),
			ID:       res.ID,
			Email:    res.Email,
			Username: res.Username,
		})
	}
}

func validateEmailorUsername(err error, log slog.Logger) (string, int) {
	switch err {
	case errorsi.ErrEmailExists:
		log.Error("Email is Exists")
		return "Email is Exists", http.StatusConflict
	case errorsi.ErrUsernameExists:
		log.Error("Username is Exists")
		return "Username is Exists", http.StatusConflict
	default:
		log.Error("DB error", sl.Err(err))
		return "Internal Error", http.StatusInternalServerError
	}

}
