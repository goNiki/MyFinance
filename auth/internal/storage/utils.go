package storage

import (
	"auth/internal/models/errorsi"
	"context"
	"fmt"
)

func (s *Storage) IsUserExistsByEmail(email string) error {
	op := "isuserexistbyemail"

	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`

	var exists bool

	err := s.DB.QueryRow(context.Background(), query, email).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if exists {
		return errorsi.ErrEmailExists
	}

	return nil
}

func (s *Storage) IsUserExistByUserName(username string) error {
	op := "isuserexistsbyusername"

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`

	var exist bool

	err := s.DB.QueryRow(context.Background(), query, username).Scan(&exist)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if exist {
		return errorsi.ErrUsernameExists
	}

	return nil
}

func (s *Storage) IsUserNotExistsByEmail(email string) error {
	op := "isuserexistbyemail"

	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`

	var exists bool

	err := s.DB.QueryRow(context.Background(), query, email).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if !exists {
		return errorsi.ErrEmailNotExists
	}

	return nil
}
