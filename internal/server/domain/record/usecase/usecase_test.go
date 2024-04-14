package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/benderr/keypass/internal/server/domain/record"
	"github.com/benderr/keypass/internal/server/domain/record/dto"
	"github.com/benderr/keypass/internal/server/domain/record/usecase"
	repomocks "github.com/benderr/keypass/internal/server/domain/record/usecase/mocks"
	mocklogger "github.com/benderr/keypass/pkg/logger/mock_logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repomocks.NewMockRecordRepo(ctrl)
	crypt := repomocks.NewMockDataCrypter(ctrl)
	uc := usecase.New(repo, crypt, mocklogger.New())

	cryptedData := []byte("123")
	userId := "user_id"
	dataType := "CREDIT"
	meta := "meta text"
	model := &dto.CredentialsRecord{
		MetaRecord: dto.MetaRecord{Meta: meta},
		Info:       dto.CredentialsInfo{Login: "test", Password: "test"},
	}

	crypt.EXPECT().Encrypt(model).Return(cryptedData, meta, nil)
	repo.EXPECT().Create(gomock.Any(), userId, cryptedData, dataType, meta).Return(true, nil)

	res, err := uc.Create(context.Background(), userId, model, dataType)

	assert.NoError(t, err, "error creating record")

	assert.Equal(t, true, res)
}

func TestUpdateRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repomocks.NewMockRecordRepo(ctrl)
	crypt := repomocks.NewMockDataCrypter(ctrl)
	uc := usecase.New(repo, crypt, mocklogger.New())

	t.Run("update success ", func(t *testing.T) {
		ID := "record id"
		cryptedData := []byte("123")
		userId := "user_id"
		meta := "meta text"
		model := &dto.CredentialsRecord{
			MetaRecord: dto.MetaRecord{Meta: meta},
			Info:       dto.CredentialsInfo{Login: "test", Password: "test"},
		}

		crypt.EXPECT().Encrypt(model).Return(cryptedData, meta, nil)
		repo.EXPECT().Update(gomock.Any(), ID, cryptedData, meta).Return(nil)
		repo.EXPECT().GetByID(gomock.Any(), ID).Return(&record.Record{UserID: userId}, nil)

		err := uc.Update(context.Background(), userId, ID, model)

		assert.NoError(t, err, "error updating record")
	})

	t.Run("update error - access denied", func(t *testing.T) {
		ID := "record id"
		userId := "user_id"
		meta := "meta text"
		model := &dto.CredentialsRecord{
			MetaRecord: dto.MetaRecord{Meta: meta},
			Info:       dto.CredentialsInfo{Login: "test", Password: "test"},
		}

		repo.EXPECT().GetByID(gomock.Any(), ID).Return(&record.Record{UserID: "other user id"}, nil)

		err := uc.Update(context.Background(), userId, ID, model)

		assert.ErrorIs(t, err, usecase.ErrAccessDenied)
	})

	t.Run("update error - encrypt error", func(t *testing.T) {
		ID := "record id"
		userId := "user_id"
		undefError := errors.New("undefined content type")
		type UndefinedRecord struct {
			testField string
		}

		model := &UndefinedRecord{
			testField: "test data",
		}

		repo.EXPECT().GetByID(gomock.Any(), ID).Return(&record.Record{UserID: userId}, nil)
		crypt.EXPECT().Encrypt(model).Return(nil, "", undefError)

		err := uc.Update(context.Background(), userId, ID, model)

		assert.ErrorIs(t, err, undefError)
	})
}

func TestGetRecords(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repomocks.NewMockRecordRepo(ctrl)
	crypt := repomocks.NewMockDataCrypter(ctrl)
	uc := usecase.New(repo, crypt, mocklogger.New())

	t.Run("get list success ", func(t *testing.T) {
		decryptedData := map[string]any{"test": 1}
		cryptedData := []byte("123")
		listItem := dto.ReadRecord{
			ID:       "ID",
			Meta:     "Meta",
			Info:     decryptedData,
			Version:  1,
			DataType: "CREDIT",
		}

		userId := "user_id"

		crypt.EXPECT().Decrypt(cryptedData).Return(decryptedData, nil)
		repoRecords := make([]record.Record, 0)
		repoRecords = append(repoRecords, record.Record{
			ID:       "ID",
			Meta:     "Meta",
			Info:     cryptedData,
			Version:  1,
			DataType: "CREDIT",
		})
		repo.EXPECT().GetByUser(gomock.Any(), userId).Return(repoRecords, nil)

		list, err := uc.GetByUser(context.Background(), userId)

		assert.NoError(t, err, "error updating record")
		assert.Equal(t, len(list), 1)
		assert.Equal(t, list[0], listItem)
	})
}
