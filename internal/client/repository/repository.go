package repository

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/benderr/keypass/internal/client/dto"
	"github.com/benderr/keypass/internal/client/session"
	"github.com/benderr/keypass/internal/secure"
	"github.com/benderr/keypass/pkg/kcrypt"
	"github.com/benderr/keypass/pkg/logger"
)

type secureRepo struct {
	path   string
	logger logger.Logger
}

func New(path string, l logger.Logger) *secureRepo {
	return &secureRepo{
		path:   path,
		logger: l,
	}
}

func (s *secureRepo) getFile(userID string) (io.ReadWriteCloser, error) {
	storageFileName := fmt.Sprintf("%s%s%s.keypass", s.path, string(os.PathSeparator), userID)

	if err := os.MkdirAll(filepath.Dir(storageFileName), 0770); err != nil {
		return nil, err
	}

	if err := s.createFileIfNotExist(storageFileName, userID); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(storageFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (s *secureRepo) createFileIfNotExist(fileName string, userID string) error {
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer file.Close()
		dec := json.NewEncoder(file)
		return dec.Encode(UserState{UserID: userID})
	}
	return nil
}

func (s *secureRepo) getUserState(userID string) (*UserState, error) {
	w, err := s.getFile(userID)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &UserState{}, nil
		}
		return nil, err
	}
	defer w.Close()

	dec := json.NewDecoder(w)
	state := new(UserState)
	err = dec.Decode(state)
	return state, err
}

func (s *secureRepo) saveUserState(state UserState) error {
	w, err := s.getFile(state.UserID)

	if err != nil {
		s.logger.Errorln("invalid writer", err)
		return err
	}

	defer w.Close()

	dec := json.NewEncoder(w)
	return dec.Encode(state)
}

func (s *secureRepo) encryptToken(token string, pin string) (string, error) {
	var buf bytes.Buffer
	enc := secure.NewEncoder(&buf)
	err := enc.Encode(token, pin)
	return hex.EncodeToString(buf.Bytes()), err
}

func (s *secureRepo) decryptToken(tokenHex string, pin string) (string, error) {
	var buf bytes.Buffer
	records, err := hex.DecodeString(tokenHex)
	if err != nil {
		return "", err
	}
	buf.Write(records)
	enc := secure.NewDecoder(&buf)
	token := ""
	err = enc.Decode(&token, pin)
	return token, err
}

func (s *secureRepo) encryptRecords(records []dto.ClientRecord, pin string) (string, error) {
	var buf bytes.Buffer
	enc := secure.NewEncoder(&buf)
	err := enc.Encode(records, pin)
	return base64.StdEncoding.EncodeToString(buf.Bytes()), err
}

func (s *secureRepo) decryptRecords(recordsHex string, pin string) ([]dto.ClientRecord, error) {
	var buf bytes.Buffer
	records, err := base64.StdEncoding.DecodeString(recordsHex)
	if err != nil {
		return nil, err
	}
	buf.Write(records)
	enc := secure.NewDecoder(&buf)
	outRecords := make([]dto.ClientRecord, 0)
	err = enc.Decode(&outRecords, pin)
	return outRecords, err
}

func (s *secureRepo) GetRecords(userID string, pin string) ([]dto.ClientRecord, error) {
	state, err := s.getUserState(userID)

	if err != nil {
		s.logger.Errorln("invalid user state", err)
		return nil, err
	}
	if len(state.Records) == 0 {
		return []dto.ClientRecord{}, nil
	}
	outRecords, err := s.decryptRecords(state.Records, pin)
	if err != nil {
		s.logger.Errorln("encryptRecords", err)
		return nil, err
	}
	return outRecords, nil
}

func (s *secureRepo) UpdateRecords(userID string, pin string, records []dto.ClientRecord) error {
	state, err := s.getUserState(userID)

	if err != nil {
		s.logger.Errorln("invalid user state", err)
		return err
	}

	inRecords, err := s.encryptRecords(records, pin)

	if err != nil {
		s.logger.Errorln("encryptRecords", err)
		return err
	}

	state.Records = inRecords
	return s.saveUserState(*state)
}

func (s *secureRepo) CreateUser(userID string, login string) (*session.UserInfo, error) {
	state, err := s.getUserState(userID)

	if err != nil {
		s.logger.Errorln("user state error", err)
		return nil, err
	}

	state.UserID = userID
	state.Login = login

	if err = s.saveUserState(*state); err != nil {
		return nil, err
	}
	return &session.UserInfo{ID: state.UserID, HashPin: state.HashPin}, err
}

func (s *secureRepo) ClearUser(userID string) error {
	state, err := s.getUserState(userID)

	if err != nil {
		s.logger.Errorln("user state error", err)
		return err
	}

	state.HashToken = ""
	return s.saveUserState(*state)
}

func (s *secureRepo) UpdateUserPin(userID string, pin string) error {
	state, err := s.getUserState(userID)

	if err != nil {
		s.logger.Errorln("invalid user state", err)
		return err
	}

	hashPin, err := kcrypt.HashString(pin)

	if err != nil {
		return err
	}

	state.HashPin = hashPin
	return s.saveUserState(*state)
}

func (s *secureRepo) UpdateUserToken(userID string, pin string, token string) error {
	state, err := s.getUserState(userID)

	if err != nil {
		s.logger.Errorln("invalid user state", err)
		return err
	}

	hashToken, err := s.encryptToken(token, pin)

	if err != nil {
		return err
	}

	state.HashToken = hashToken
	return s.saveUserState(*state)
}

func (s *secureRepo) LoadLastUser() (*session.UserInfo, error) {
	dir := fmt.Sprintf("%s%s", s.path, string(os.PathSeparator))

	if err := os.MkdirAll(dir, 0770); err != nil {
		return nil, err
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, nil
	}

	lastUpdated := time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC)
	userID := ""
	for _, f := range files {
		info, err := f.Info()
		if err != nil {
			continue
		}
		if info.ModTime().After(lastUpdated) {
			lastUpdated = info.ModTime()
			userID = strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
		}
	}

	if len(userID) == 0 {
		return nil, nil
	}

	state, err := s.getUserState(userID)

	if err != nil {
		return nil, err
	}

	if len(state.UserID) > 0 && len(state.HashToken) > 0 {
		return &session.UserInfo{ID: state.UserID, HashPin: state.HashPin, HashToken: state.HashToken}, nil
	}
	return nil, nil
}

func (s *secureRepo) CheckUserPin(userID string, pin string) (bool, string, error) {
	state, err := s.getUserState(userID)

	if err != nil {
		s.logger.Errorln("invalid user state", err)
		return false, "", err
	}

	if len(state.HashPin) == 0 {
		s.logger.Errorln("pin didn't set", err)
		return false, "", errors.New("pin didn't set")
	}

	if !kcrypt.CheckString(pin, state.HashPin) {
		return false, "", nil
	}

	token := ""
	if len(state.HashToken) > 0 {
		token, err = s.decryptToken(state.HashToken, pin)
		if err != nil {
			return false, "", err
		}
	}

	return true, token, nil
}
