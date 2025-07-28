package storage

import (
	"auth/internal/config"
	"auth/internal/handlers/register"
	"auth/internal/models/users"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStorage interface {
	RegisterUser(ctx context.Context, email string, username string, PassHash []byte, createat time.Time) (*register.Response, error)
	GetUserByEmail(ctx context.Context, email string) (*users.Users, error)
	IsUserExistsByEmail(email string) error
	IsUserExistByUserName(username string) error
	IsUserNotExistsByEmail(email string) error
}

type Storage struct {
	DB *pgxpool.Pool
}

func New(cfg *config.DBConfig) (*Storage, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBname, cfg.SSLmode)

	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable  to connect to database %w", err)
	}
	if err := db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping to database %w ", err)
	}

	return &Storage{DB: db}, nil
}

// RegisterUser registers a new user in the database.
func (s *Storage) RegisterUser(ctx context.Context, email string, username string, PassHash []byte, createat time.Time) (*register.Response, error) {
	const op = "storage.registeruser"

	query := `INSERT INTO users (email, username, password, createat) VALUES ($1, $2, $3, $4) RETURNING id`

	var userID int64

	err := s.DB.QueryRow(ctx, query, email, username, PassHash, createat).Scan(&userID)
	if err != nil {
		//TODO сделать обработчик ошибки нет пользователя
		return &register.Response{}, fmt.Errorf("%s : %w", op, err)
	}

	return &register.Response{ID: userID, Email: email, Username: username}, nil

}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*users.Users, error) {

	const op = "storage.getuserbyemail"

	query := `SELECT id, email, password FROM users WHERE email = $1`

	var user users.Users

	err := s.DB.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return &users.Users{}, fmt.Errorf("%s: %w", op, err)
	}

	return &users.Users{ID: user.ID, Email: user.Email, Password: user.Password}, nil
}
