package storage

import (
	"auth/internal/config"
	"auth/internal/handlers/register"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
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

	return &Storage{db: db}, nil
}

// RegisterUser registers a new user in the database.
func (s *Storage) RegisterUser(ctx context.Context, email string, username string, PassHash []byte, createat time.Time) (*register.Response, error) {
	const op = "storage.registeruser"

	query := `INSERT INTO users (email, username, password, createat) VALUES ($1, $2, $3, $4) RETURNING id`

	var userID int64

	err := s.db.QueryRow(ctx, query, email, username, PassHash, createat).Scan(&userID)
	if err != nil {
		//TODO сделать обработчик ошибки нет пользователя
		return &register.Response{}, fmt.Errorf("%s : %w", op, err)
	}

	return &register.Response{ID: userID, Email: email, Username: username}, nil

}
