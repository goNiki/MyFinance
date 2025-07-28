package users

import "time"

type Users struct {
	ID       int64
	Email    string
	Username string
	Password string
	CreateAt time.Time
}
