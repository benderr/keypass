package record

import (
	"errors"
	"time"
)

type DataType = string

type Record struct {
	ID        string    `json:"id" validate:"required"`
	Meta      string    `json:"meta"`
	Info      []byte    `json:"info" validate:"required"`
	Version   int       `json:"version,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
	DataType  DataType  `json:"data_type"`
	UserID    string    `json:"-"`
}

var (
	ErrAlreadyExist = errors.New("already exist")
	ErrNotFound     = errors.New("not found")
)

const (
	CREDENTIALS DataType = "CREDENTIALS"
	TEXT        DataType = "TEXT"
	BINARY      DataType = "BINARY"
	CREDIT      DataType = "CREDIT"
)
