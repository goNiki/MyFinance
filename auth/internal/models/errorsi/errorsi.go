package errorsi

import "errors"

var (
	ErrEmailExists    = errors.New("email exists")
	ErrEmailNotExists = errors.New("email not exists")
	ErrUsernameExists = errors.New("username exists")
	ErrInvaliPassword = errors.New("invalid password")
)
