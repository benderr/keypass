package usecase

import (
	"context"
	"errors"

	"github.com/benderr/keypass/internal/server/domain/record"
	"github.com/benderr/keypass/internal/server/domain/record/dto"
	"github.com/benderr/keypass/pkg/logger"
)

type recordlUsecase struct {
	recordRepo RecordRepo
	logger     logger.Logger
	crypt      DataCrypter
}

var (
	ErrNotFound     = errors.New("not found")
	ErrAccessDenied = errors.New("access denied")
)

// New return instance with methods for http handlers
func New(op RecordRepo, crypt DataCrypter, l logger.Logger) *recordlUsecase {
	return &recordlUsecase{
		recordRepo: op,
		logger:     l,
		crypt:      crypt,
	}
}

func (r *recordlUsecase) Update(ctx context.Context, userID string, ID string, inModel any) error {
	record, err := r.recordRepo.GetByID(ctx, ID)
	if err != nil {
		return err
	}
	if record.UserID != userID {
		return ErrAccessDenied
	}

	content, meta, err := r.crypt.Encrypt(inModel)
	if err != nil {
		return err
	}
	return r.recordRepo.Update(ctx, ID, content, meta)
}

func (r *recordlUsecase) Create(ctx context.Context, userID string, inModel any, dataType record.DataType) (bool, error) {
	content, meta, err := r.crypt.Encrypt(inModel)
	if err != nil {
		return false, err
	}
	return r.recordRepo.Create(ctx, userID, content, dataType, meta)
}

func (r *recordlUsecase) GetByUser(ctx context.Context, userid string) ([]dto.ReadRecord, error) {
	list, err := r.recordRepo.GetByUser(ctx, userid)
	if err != nil {
		return nil, err
	}
	res := make([]dto.ReadRecord, 0)
	for _, rec := range list {

		info, err := r.crypt.Decrypt(rec.Info)

		if err != nil {
			return nil, err
		}

		res = append(res, dto.ReadRecord{
			ID:        rec.ID,
			Meta:      rec.Meta,
			Info:      info,
			Version:   rec.Version,
			UpdatedAt: rec.UpdatedAt,
			DataType:  rec.DataType,
		})
	}
	return res, nil
}

func (r *recordlUsecase) Delete(ctx context.Context, userID string, ID string) error {
	record, err := r.recordRepo.GetByID(ctx, ID)
	if err != nil {
		return err
	}
	if record.UserID != userID {
		return ErrAccessDenied
	}

	return r.recordRepo.Delete(ctx, ID)
}
