package datacrypt

import (
	"bytes"
	"errors"

	"github.com/benderr/keypass/internal/secure"
	"github.com/benderr/keypass/internal/server/domain/record/dto"
	"github.com/benderr/keypass/pkg/logger"
)

type dataCrypt struct {
	secret string
	logger logger.Logger
}

var ErrUndefinedType = errors.New("undefined content type")

// New return structure for encrypt and decrypt records
func New(secret string, logger logger.Logger) *dataCrypt {
	return &dataCrypt{
		secret: secret,
		logger: logger,
	}
}

// Encrypt create encrypted record
func (r *dataCrypt) Encrypt(inModel any) ([]byte, string, error) {
	var buf bytes.Buffer
	var meta string
	content := make(map[string]any, 0)
	switch v := inModel.(type) {
	case *dto.CredentialsRecord:
		content["login"] = v.Info.Login
		content["password"] = v.Info.Password
		meta = v.Meta
	case *dto.CreditCardRecord:
		content["number"] = v.Info.Number
		content["cvv"] = v.Info.CVV
		content["expire"] = v.Info.Expire
		meta = v.Meta
	case *dto.TextRecord:
		content["text"] = v.Info.Text
		meta = v.Meta
	case *dto.BinaryRecord:
		content["binary"] = v.Data
		content["filePath"] = v.FilePath
		meta = v.Meta
	default:
		r.logger.Infoln(v)
		return nil, "", ErrUndefinedType
	}

	enc := secure.NewEncoder(&buf)
	err := enc.Encode(content, r.secret)
	if err != nil {
		return nil, "", err
	}
	return buf.Bytes(), meta, nil
}

// Decrypt try to decrypt record to map
func (r *dataCrypt) Decrypt(model []byte) (map[string]any, error) {
	var buf bytes.Buffer
	_, err := buf.Write(model)

	if err != nil {
		return nil, err
	}

	dec := secure.NewDecoder(&buf)

	content := make(map[string]any, 0)

	err = dec.Decode(&content, r.secret)
	return content, err
}
