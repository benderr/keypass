package datacrypt_test

import (
	"testing"

	"github.com/benderr/keypass/internal/server/domain/record/datacrypt"
	"github.com/benderr/keypass/internal/server/domain/record/dto"
	mocklogger "github.com/benderr/keypass/pkg/logger/mock_logger"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	dc := datacrypt.New("123", mocklogger.New())

	type want struct {
		result map[string]any
		meta   string
	}

	tests := []struct {
		dataType string
		payload  interface{}
		want     want
	}{
		{
			dataType: "CREDENTIALS",
			payload: &dto.CredentialsRecord{
				MetaRecord: dto.MetaRecord{Meta: "Test meta"},
				Info:       dto.CredentialsInfo{Login: "test", Password: "test"},
			},
			want: want{
				meta:   "Test meta",
				result: map[string]any{"login": "test", "password": "test"},
			},
		},
		{
			dataType: "TEXT",
			payload: &dto.TextRecord{
				MetaRecord: dto.MetaRecord{Meta: "Test meta for text"},
				Info:       dto.TextInfo{Text: "secret text"},
			},
			want: want{
				meta:   "Test meta for text",
				result: map[string]any{"text": "secret text"},
			},
		},
		{
			dataType: "CREDIT",
			payload: &dto.CreditCardRecord{
				MetaRecord: dto.MetaRecord{Meta: "Test meta for credit card"},
				Info:       dto.CreditCardInfo{Number: "123", CVV: "123", Expire: "12/24"},
			},
			want: want{
				meta:   "Test meta for credit card",
				result: map[string]any{"number": "123", "cvv": "123", "expire": "12/24"},
			},
		},
		{
			dataType: "BINARY",
			payload: &dto.BinaryRecord{
				MetaRecord: dto.MetaRecord{Meta: "Test meta for binary"},
				Data:       []byte("123"),
				FilePath:   "c:/test/test.txt",
			},
			want: want{
				meta:   "Test meta for binary",
				result: map[string]any{"binary": []byte("123"), "filePath": "c:/test/test.txt"},
			},
		},
	}

	for _, test := range tests {
		t.Run(" Encrypt and decrypt for "+test.dataType, func(t *testing.T) {
			res, meta, err := dc.Encrypt(test.payload)
			assert.NoError(t, err, "encrypting error")

			descryptResult, err := dc.Decrypt(res)
			assert.NoError(t, err, "decrypt error")

			assert.Equal(t, test.want.meta, meta)
			assert.Equal(t, test.want.result, descryptResult)
		})
	}
}

func TestEncryptFailed(t *testing.T) {
	dc := datacrypt.New("123", mocklogger.New())

	type UndefRecord struct {
		Test string
	}

	_, _, err := dc.Encrypt(&UndefRecord{
		Test: "123",
	})

	assert.ErrorIs(t, datacrypt.ErrUndefinedType, err)
}
