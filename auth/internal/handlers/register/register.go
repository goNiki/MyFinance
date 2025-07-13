package register

import (
	resp "auth/internal/utils/api/response"
	"auth/internal/utils/logger/sl"
	"context"
	"fmt"
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

type Register interface {
	RegisterUser(ctx context.Context, email string, username string, PassHash []byte, createat time.Time) (*Response, error)
}

func New(log *slog.Logger, register Register) http.HandlerFunc {
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
			log.Error("invalid request", sl.Err(err))
			//TODO доделать обработчик ошибок после валидации.
			fmt.Println(validateErr)
			http.Error(w, "validator err", http.StatusBadRequest)
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
			//TODO переделать блок обработки ошибок
			log.Error("error save BD")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		log.Info("user register", slog.Int64("ID", int64(res.ID)), slog.String("email", res.Email))

		render.JSON(w, r, Response{
			Responce: resp.Ok(),
			ID:       res.ID,
			Email:    res.Email,
			Username: res.Username,
		})
	}

}
