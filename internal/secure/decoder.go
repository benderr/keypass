package secure

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"io"
	"reflect"
)

type cryptoDecoder struct {
	r io.Reader
}

func NewDecoder(r io.Reader) *cryptoDecoder {
	dec := new(cryptoDecoder)
	// We use the ability to read bytes as a plausible surrogate for buffering.
	if _, ok := r.(io.ByteReader); !ok {
		r = bufio.NewReader(r)
	}
	dec.r = r
	return dec
}

var ErrIsEmpty = errors.New("decode empty value")
var ErrNotPointer = errors.New("attempt to decode into a non-pointer")

func (c *cryptoDecoder) Decode(e any, masterKey string) error {
	if e == nil {
		return ErrIsEmpty
	}
	value := reflect.ValueOf(e)
	if value.Type().Kind() != reflect.Pointer {
		return ErrNotPointer
	}

	buf := &bytes.Buffer{}
	teeReader := io.TeeReader(c.r, buf)
	secureContent, err := io.ReadAll(teeReader)

	if err != nil {
		return err
	}

	unsecContent, err := decryptMessage(secureContent, masterKey)

	if err != nil {
		return err
	}

	var decBuf bytes.Buffer
	decBuf.Write(unsecContent)
	dec := gob.NewDecoder(&decBuf)

	return dec.Decode(e)

}

func decryptMessage(message []byte, password string) ([]byte, error) {
	key := sha256.Sum256([]byte(password))

	// создаем aesblock и aesgcm
	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	// создаём вектор инициализации
	nonce := key[len(key)-aesgcm.NonceSize():]

	// расшифровываем
	decrypted, err := aesgcm.Open(nil, nonce, message, nil)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}
