package logic_test

import (
	"testing"

	"github.com/benderr/keypass/internal/client/dto"
	"github.com/benderr/keypass/internal/client/logic"
	"github.com/benderr/keypass/internal/client/logic/logicmocks"
	"github.com/benderr/keypass/internal/client/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLoadUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queryClient := logicmocks.NewMockIQueryClient(ctrl)
	repo := logicmocks.NewMockSecureRepository(ctrl)

	l := logic.New(queryClient, repo)

	repo.EXPECT().LoadLastUser().Return(&session.UserInfo{ID: "test", HashPin: "hash"}, nil)
	err := l.LoadUser()
	require.NoError(t, err)
	state := l.GetSessionState()
	assert.Equal(t, session.Suspended, state)
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queryClient := logicmocks.NewMockIQueryClient(ctrl)
	repo := logicmocks.NewMockSecureRepository(ctrl)

	l := logic.New(queryClient, repo)
	id := "user_id"
	login := "login"
	pass := "pass"

	queryClient.EXPECT().Login(login, pass).Return(dto.User{ID: id, Login: login}, nil)
	repo.EXPECT().CreateUser(id, login).Return(&session.UserInfo{ID: id}, nil)

	err := l.Login(login, pass)
	require.NoError(t, err)
	state := l.GetSessionState()
	assert.Equal(t, session.NeedPin, state)
}

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queryClient := logicmocks.NewMockIQueryClient(ctrl)
	repo := logicmocks.NewMockSecureRepository(ctrl)

	l := logic.New(queryClient, repo)
	id := "user_id"
	login := "login"
	pass := "pass"

	queryClient.EXPECT().Register(login, pass).Return(dto.User{ID: id, Login: login}, nil)
	repo.EXPECT().CreateUser(id, login).Return(&session.UserInfo{ID: id}, nil)

	err := l.Register(login, pass)
	require.NoError(t, err)
	state := l.GetSessionState()
	assert.Equal(t, session.NeedPin, state)
}

func TestLogout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queryClient := logicmocks.NewMockIQueryClient(ctrl)
	repo := logicmocks.NewMockSecureRepository(ctrl)

	l := logic.New(queryClient, repo)
	id := "user_id"

	repo.EXPECT().LoadLastUser().Return(&session.UserInfo{ID: id}, nil)
	repo.EXPECT().ClearUser(id).Return(nil)
	l.LoadUser()

	err := l.Logout()
	require.NoError(t, err)
	state := l.GetSessionState()
	assert.Equal(t, session.NoSession, state)
}

func TestSuspendSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queryClient := logicmocks.NewMockIQueryClient(ctrl)
	repo := logicmocks.NewMockSecureRepository(ctrl)

	l := logic.New(queryClient, repo)
	id := "user_id"
	login := "test_login"
	pass := "test_pass"
	pin := "test_pin"
	token := "token"

	queryClient.EXPECT().Login(login, pass).Return(dto.User{ID: id, Login: login, Token: token}, nil)
	repo.EXPECT().CreateUser(id, login).Return(&session.UserInfo{ID: id, HashPin: ""}, nil)
	err := l.Login(login, pass)
	require.NoError(t, err)

	repo.EXPECT().UpdateUserPin(id, pin).Return(nil)
	repo.EXPECT().UpdateUserToken(id, pin, token).Return(nil)
	err = l.SetPin(pin)
	require.NoError(t, err)

	stateActive := l.GetSessionState()
	assert.Equal(t, session.Active, stateActive)

	err = l.SuspendSession()
	require.NoError(t, err)
	stateSuspended := l.GetSessionState()
	assert.Equal(t, session.Suspended, stateSuspended)
}

func TestCheckPin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queryClient := logicmocks.NewMockIQueryClient(ctrl)
	repo := logicmocks.NewMockSecureRepository(ctrl)

	l := logic.New(queryClient, repo)
	id := "user_id"
	login := "test_login"
	pass := "test_pass"
	pin := "test_pin"
	token := "token"

	queryClient.EXPECT().Login(login, pass).Return(dto.User{ID: id, Login: login, Token: token}, nil)
	repo.EXPECT().CreateUser(id, login).Return(&session.UserInfo{ID: id, HashPin: "123"}, nil)
	err := l.Login(login, pass)
	require.NoError(t, err)

	repo.EXPECT().CheckUserPin(id, pin).Return(true, token, nil)

	valid, err := l.CheckPin(pin)
	require.NoError(t, err)
	assert.Equal(t, true, valid)

	stateActive := l.GetSessionState()
	assert.Equal(t, session.Active, stateActive)
}

func TestSyncRecords(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queryClient := logicmocks.NewMockIQueryClient(ctrl)
	repo := logicmocks.NewMockSecureRepository(ctrl)

	l := logic.New(queryClient, repo)
	id := "user_id"
	login := "test_login"
	pass := "test_pass"
	pin := "test_pin"
	token := "token"

	queryClient.EXPECT().Login(login, pass).Return(dto.User{ID: id, Login: login, Token: token}, nil)
	repo.EXPECT().CreateUser(id, login).Return(&session.UserInfo{ID: id, HashPin: "123"}, nil)
	l.Login(login, pass)

	repo.EXPECT().CheckUserPin(id, pin).Return(true, token, nil)
	l.CheckPin(pin)

	records := []dto.ClientRecord{{ID: "1"}}
	queryClient.EXPECT().GetRecords(token).Return(records, nil)
	repo.EXPECT().UpdateRecords(id, pin, records).Return(nil)
	err := l.SyncRecords()
	assert.NoError(t, err)
}

func TestAddRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queryClient := logicmocks.NewMockIQueryClient(ctrl)
	repo := logicmocks.NewMockSecureRepository(ctrl)

	l := logic.New(queryClient, repo)
	id := "user_id"
	login := "test_login"
	pass := "test_pass"
	pin := "test_pin"
	token := "token"

	queryClient.EXPECT().Login(login, pass).Return(dto.User{ID: id, Login: login, Token: token}, nil)
	repo.EXPECT().CreateUser(id, login).Return(&session.UserInfo{ID: id, HashPin: "123"}, nil)
	l.Login(login, pass)

	repo.EXPECT().CheckUserPin(id, pin).Return(true, token, nil)
	l.CheckPin(pin)

	record := dto.ServerRecord{
		ID:       "ID",
		Meta:     "Meta",
		DataType: "CREDIT",
	}
	queryClient.EXPECT().AddRecord(token, record).Return(nil)
	queryClient.EXPECT().GetRecords(token).Return([]dto.ClientRecord{{ID: "1"}}, nil)
	repo.EXPECT().UpdateRecords(id, pin, []dto.ClientRecord{{ID: "1"}}).Return(nil)

	err := l.AddRecord(record)
	assert.NoError(t, err)
}

func TestGetRecords(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queryClient := logicmocks.NewMockIQueryClient(ctrl)
	repo := logicmocks.NewMockSecureRepository(ctrl)

	l := logic.New(queryClient, repo)
	id := "user_id"
	login := "test_login"
	pass := "test_pass"
	pin := "test_pin"
	token := "token"

	queryClient.EXPECT().Login(login, pass).Return(dto.User{ID: id, Login: login, Token: token}, nil)
	repo.EXPECT().CreateUser(id, login).Return(&session.UserInfo{ID: id, HashPin: "123"}, nil)
	l.Login(login, pass)

	repo.EXPECT().CheckUserPin(id, pin).Return(true, token, nil)
	l.CheckPin(pin)

	records := []dto.ClientRecord{{ID: "1"}}
	repo.EXPECT().GetRecords(id, pin).Return(records, nil)

	result, err := l.GetRecords()
	assert.Equal(t, records, result)
	assert.NoError(t, err)
}
