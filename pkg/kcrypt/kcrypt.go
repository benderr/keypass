package kcrypt

import (
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func HashBytes(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return nil, err
	}
	return bytes, err
}

func CheckBytes(password, pwd []byte) bool {
	err := bcrypt.CompareHashAndPassword(pwd, password)
	return err == nil
}

func HashString(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), err
}

func CheckString(password, hashedPwd string) bool {
	passwordBytes := []byte(password)
	hashedPwdBytes, err := hex.DecodeString(hashedPwd)
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword(hashedPwdBytes, passwordBytes)
	return err == nil
}
