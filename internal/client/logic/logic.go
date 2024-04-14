package logic

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/benderr/keypass/internal/client/dto"
	"github.com/benderr/keypass/internal/client/session"
	recordform "github.com/benderr/keypass/pkg/client/component/record_form"
)

var (
	ErrProfileNotFound = errors.New("profile doesn't exist")
	ErrNoSession       = errors.New("session not found")
)

type appLogic struct {
	query         IQueryClient
	repo          SecureRepository
	sessionState  session.State
	sessionPin    string
	sessionToken  string
	currentUserID string
}

type IQueryClient interface {
	Login(login string, pass string) (dto.User, error)
	Register(login string, pass string) (dto.User, error)
	GetRecords(token string) ([]dto.ClientRecord, error)
	UpdateRecord(token string, record dto.ServerRecord) error
	AddRecord(token string, record dto.ServerRecord) error
	AddRecordFile(token string, record dto.ServerRecord) error
	DeleteRecord(token string, ID string) error
}

type SecureRepository interface {
	GetRecords(userID string, pin string) ([]dto.ClientRecord, error)
	UpdateRecords(userID string, pin string, records []dto.ClientRecord) error

	CreateUser(userID string, login string) (*session.UserInfo, error)
	ClearUser(userID string) error
	LoadLastUser() (*session.UserInfo, error)
	UpdateUserPin(userID string, pin string) error
	UpdateUserToken(userID string, pin string, token string) error
	CheckUserPin(userID string, pin string) (bool, string, error)
}

// New return appLogic which contain business logic for TUI
func New(q IQueryClient, r SecureRepository) *appLogic {
	return &appLogic{
		query:        q,
		repo:         r,
		sessionState: session.NoSession,
	}
}

func (a *appLogic) destroySession() error {
	if err := a.repo.ClearUser(a.currentUserID); err != nil {
		return err
	}
	a.currentUserID = ""
	a.sessionToken = ""
	a.sessionPin = ""
	a.sessionState = session.NoSession
	return nil
}

func (a *appLogic) createSession(user *dto.User) error {
	u, err := a.repo.CreateUser(user.ID, user.Login)

	if err != nil {
		return err
	}

	if len(u.HashPin) == 0 {
		a.sessionState = session.NeedPin
	} else {
		a.sessionState = session.Suspended
	}

	a.currentUserID = user.ID
	a.sessionToken = user.Token
	return nil
}

func (a *appLogic) activateSession(pin string) {
	a.sessionPin = pin
	a.sessionState = session.Active
}

// LoadUser load last active user from cache
func (a *appLogic) LoadUser() error {
	u, err := a.repo.LoadLastUser()
	if err != nil {
		return err
	}
	if u != nil {
		if len(u.HashPin) == 0 {
			a.sessionState = session.NeedPin
		} else {
			a.sessionState = session.Suspended
		}
		a.currentUserID = u.ID
	} else {
		a.sessionState = session.NoSession
	}
	return nil
}

// Login fetch user info from server by login/pass
func (a *appLogic) Login(login string, pass string) error {
	u, err := a.query.Login(login, pass)
	if err != nil {
		return err
	}
	return a.createSession(&u)
}

// Login create new user on server
func (a *appLogic) Register(login string, pass string) error {
	u, err := a.query.Register(login, pass)
	if err != nil {
		return err
	}
	return a.createSession(&u)
}

// Logout destroy active user session
func (a *appLogic) Logout() error {
	return a.destroySession()
}

// Logout destroy active user session
func (a *appLogic) GetRecords() ([]dto.ClientRecord, error) {
	if a.sessionState != session.Active {
		return []dto.ClientRecord{}, nil
	}

	return a.repo.GetRecords(a.currentUserID, a.sessionPin)
}

// DeleteRecord delete user record by id and sync records
func (a *appLogic) DeleteRecord(ID string) error {
	if err := a.query.DeleteRecord(a.sessionToken, ID); err != nil {
		return err
	}
	return a.SyncRecords()
}

// UpdateRecord modified exist record and sync records
func (a *appLogic) UpdateRecord(record dto.ServerRecord) error {
	if err := a.query.UpdateRecord(a.sessionToken, record); err != nil {
		return err
	}
	return a.SyncRecords()
}

// AddRecord add new record to server and sync records
func (a *appLogic) AddRecord(record dto.ServerRecord) error {
	if record.DataType == recordform.BINARY {
		if err := a.query.AddRecordFile(a.sessionToken, record); err != nil {
			return err
		}
	} else {
		if err := a.query.AddRecord(a.sessionToken, record); err != nil {
			return err
		}
	}

	return a.SyncRecords()
}

// SyncRecords fetch records from server and save to user cache
func (a *appLogic) SyncRecords() error {
	if a.sessionState != session.Active {
		return nil
	}

	records, err := a.query.GetRecords(a.sessionToken)

	if err != nil {
		return err
	}

	return a.repo.UpdateRecords(a.currentUserID, a.sessionPin, records)
}

// GetSessionState return current user session state
func (a *appLogic) GetSessionState() session.State {
	return a.sessionState
}

// SuspendSession deactivate user session state
func (a *appLogic) SuspendSession() error {
	if a.sessionState == session.Active {
		a.sessionPin = ""
		a.sessionState = session.Suspended
	}
	return nil
}

// CheckPin validate user pin and activate session if valid
func (a *appLogic) CheckPin(pin string) (bool, error) {
	if a.sessionState != session.Suspended && a.sessionState != session.Active {
		return false, ErrNoSession
	}

	if len(a.currentUserID) == 0 {
		return false, ErrProfileNotFound
	}

	valid, token, err := a.repo.CheckUserPin(a.currentUserID, pin)

	if err != nil {
		return false, err
	}

	if valid {
		a.activateSession(pin)

		if len(token) == 0 && len(a.sessionToken) > 0 {
			if err := a.repo.UpdateUserToken(a.currentUserID, pin, a.sessionToken); err != nil {
				return false, err
			}
		} else {
			a.sessionToken = token
		}
	}

	return valid, nil
}

// SetPin set new pin for current user
func (a *appLogic) SetPin(pin string) error {
	if len(a.currentUserID) == 0 {
		return ErrProfileNotFound
	}

	if err := a.repo.UpdateUserPin(a.currentUserID, pin); err != nil {
		return err
	}

	if err := a.repo.UpdateUserToken(a.currentUserID, pin, a.sessionToken); err != nil {
		return err
	}

	a.activateSession(pin)
	return nil
}

// SaveBinaryFile save existed record byte array to cache dir
func (a *appLogic) SaveBinaryFile(r *dto.ClientRecord) (string, error) {
	b, ok := r.Info["binary"].(string)

	if !ok {
		return "", errors.New("file corrupted")
	}

	sDec, err := base64.StdEncoding.DecodeString(b)

	if err != nil {
		return "", err
	}

	fileName := time.Now().Format(time.RFC3339) + ".temp"
	if f, ok := r.Info["filePath"].(string); ok {
		_, fName := filepath.Split(f)
		if len(fName) > 0 {
			fileName = fName
		}
	}

	dir, _ := os.UserCacheDir()
	filePath := fmt.Sprintf("%s%s%s", dir, string(os.PathSeparator), fileName)
	f, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	f.Write(sDec)
	defer f.Close()
	return filePath, nil
}
