package usecase

import (
	"context"

	"github.com/benderr/keypass/internal/server/domain/user"
	"github.com/benderr/keypass/pkg/kcrypt"
	"github.com/benderr/keypass/pkg/logger"
)

type userRepo interface {
	GetUserByLogin(ctx context.Context, login string) (*user.User, error)
	AddUser(ctx context.Context, login string, passhash []byte) (*user.User, error)
}

type userUsecase struct {
	repo   userRepo
	logger logger.Logger
}

func New(repo userRepo, logger logger.Logger) *userUsecase {
	return &userUsecase{repo: repo, logger: logger}
}

func (u *userUsecase) Login(ctx context.Context, login, password string) (*user.User, error) {
	usr, err := u.repo.GetUserByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	if usr == nil {
		return nil, user.ErrNotFound
	}

	if !kcrypt.CheckBytes([]byte(password), usr.Password) {
		return nil, user.ErrBadPass
	}

	return usr, nil
}

func (u *userUsecase) Register(ctx context.Context, login, password string) (*user.User, error) {
	passhash, err := kcrypt.HashBytes(password)

	if err != nil {
		u.logger.Errorln("HASH ERROR", err)
		return nil, err
	}

	u.logger.Infoln("HASH PASSWORD", passhash)
	return u.repo.AddUser(ctx, login, passhash)
}
