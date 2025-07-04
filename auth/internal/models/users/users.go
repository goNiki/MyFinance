package users

import "time"

type Users struct {
	ID       int
	Email    string
	Username string
	Password string
	CreateAt time.Time
}
