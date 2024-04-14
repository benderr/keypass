package app

import (
	"context"
	"net/http"

	"github.com/benderr/keypass/internal/server/config"

	userDelivery "github.com/benderr/keypass/internal/server/domain/user/delivery"
	userRepository "github.com/benderr/keypass/internal/server/domain/user/repository"
	userUsecase "github.com/benderr/keypass/internal/server/domain/user/usecase"

	"github.com/benderr/keypass/internal/server/domain/record/datacrypt"
	recordDelivery "github.com/benderr/keypass/internal/server/domain/record/delivery"
	recordRepository "github.com/benderr/keypass/internal/server/domain/record/repository"
	recordUsecase "github.com/benderr/keypass/internal/server/domain/record/usecase"

	"github.com/benderr/keypass/internal/server/migration"
	"github.com/benderr/keypass/internal/server/session"
	"github.com/benderr/keypass/pkg/logger"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func Run(ctx context.Context, conf *config.Config) {
	logger, sync := logger.New()
	defer sync()

	db := migration.MustLoad(ctx, conf, logger)

	sessionManager := session.New(conf.SecretKey)

	userRepo := userRepository.New(db, logger)
	recordRepo := recordRepository.New(db, logger)

	userUse := userUsecase.New(userRepo, logger)
	crypt := datacrypt.New(conf.RecordSecretKey, logger)
	recordUse := recordUsecase.New(recordRepo, crypt, logger)

	e := echo.New()
	validate := validator.New()

	e.Validator = &CustomValidator{validator: validate}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())

	publicGroup := e.Group("")

	privateGroup := e.Group("", echojwt.WithConfig(echojwt.Config{
		SigningKey:    []byte(conf.SecretKey),
		NewClaimsFunc: func(c echo.Context) jwt.Claims { return new(session.UserClaims) },
	}))

	userDelivery.NewUserHandlers(publicGroup, userUse, sessionManager, logger)
	recordDelivery.NewRecordHandlers(privateGroup, recordUse, sessionManager, logger)

	e.Logger.Fatal(e.StartTLS(string(conf.Server), conf.PublicKey, conf.PrivateKey))
}
