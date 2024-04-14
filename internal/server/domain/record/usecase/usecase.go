package usecase

import (
	"context"

	"github.com/benderr/keypass/internal/server/domain/record"
)

type RecordRepo interface {
	Update(ctx context.Context, ID string, info []byte, meta string) error
	Delete(ctx context.Context, ID string) error
	Create(ctx context.Context, userID string, info []byte, dataType record.DataType, meta string) (bool, error)
	GetByUser(ctx context.Context, userID string) ([]record.Record, error)
	GetByID(ctx context.Context, ID string) (*record.Record, error)
}

type DataCrypter interface {
	Encrypt(model any) ([]byte, string, error)
	Decrypt(model []byte) (map[string]any, error)
}
