package dto

import "time"

type MetaRecord struct {
	Meta string `json:"meta" validate:"required"`
}

type CredentialsInfo struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type CredentialsRecord struct {
	MetaRecord
	Info CredentialsInfo `json:"info" validate:"required"`
}

type TextInfo struct {
	Text string `json:"text" validate:"required"`
}
type TextRecord struct {
	MetaRecord
	Info TextInfo `json:"info" validate:"required"`
}

type BinaryRecord struct {
	MetaRecord
	Data     []byte `json:"data" validate:"required"`
	FilePath string
}

type CreditCardInfo struct {
	Number string `json:"number" validate:"required"`
	CVV    string `json:"cvv" validate:"required"`
	Expire string `json:"expire" validate:"required"`
}
type CreditCardRecord struct {
	MetaRecord
	Info CreditCardInfo `json:"info" validate:"required"`
}

type ReadRecord struct {
	ID        string         `json:"id"`
	Meta      string         `json:"meta"`
	Info      map[string]any `json:"info"`
	Version   int            `json:"version"`
	UpdatedAt time.Time      `json:"updated_at"`
	DataType  string         `json:"data_type"`
}
