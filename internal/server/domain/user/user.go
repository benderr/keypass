package user

import (
	"errors"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Login     string    `json:"login" validate:"required"`
	Password  []byte    `json:"password" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	ErrBadPass    = errors.New("bad pass")
	ErrNotFound   = errors.New("user not found")
	ErrLoginExist = errors.New("login already exist")
)
