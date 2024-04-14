package secure

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/gob"
	"io"
)

type cryptoEncoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *cryptoEncoder {
	return &cryptoEncoder{
		w: w,
	}
}

func (c *cryptoEncoder) Encode(e any, masterKey string) error {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)

	err := enc.Encode(e)

	if err != nil {
		return err
	}

	res, err := encryptMessage(network.Bytes(), masterKey)

	if err != nil {
		return err
	}

	_, err2 := c.w.Write(res)
	return err2
}

func encryptMessage(src []byte, password string) ([]byte, error) {
	key := sha256.Sum256([]byte(password))

	// NewCipher создает и возвращает новый cipher.Block.
	// Ключевым аргументом должен быть ключ AES, 16, 24 или 32 байта
	// для выбора AES-128, AES-192 или AES-256.
	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	// NewGCM возвращает заданный 128-битный блочный шифр
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	// создаём вектор инициализации
	nonce := key[len(key)-aesgcm.NonceSize():]

	dst := aesgcm.Seal(nil, nonce, src, nil) // зашифровываем

	return dst, nil
}
