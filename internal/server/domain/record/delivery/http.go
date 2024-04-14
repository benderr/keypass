package delivery

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/benderr/keypass/internal/server/domain/record"
	"github.com/benderr/keypass/internal/server/domain/record/dto"
	"github.com/benderr/keypass/pkg/httputils"
	"github.com/benderr/keypass/pkg/logger"
	"github.com/labstack/echo/v4"
)

type RecordUsecase interface {
	Create(ctx context.Context, userID string, dto any, dataType record.DataType) (bool, error)
	Update(ctx context.Context, userID string, ID string, inModel any) error
	Delete(ctx context.Context, userID string, ID string) error
	GetByUser(ctx context.Context, userID string) ([]dto.ReadRecord, error)
}

type SessionManager interface {
	GetUserID(c echo.Context) (string, error)
}

type recordsHandler struct {
	session SessionManager
	logger  logger.Logger
	RecordUsecase
}

func NewRecordHandlers(group *echo.Group, or RecordUsecase, session SessionManager, logger logger.Logger) {
	h := &recordsHandler{
		RecordUsecase: or,
		session:       session,
		logger:        logger,
	}

	g := group.Group("/api")

	g.GET("/records", h.GetRecordsHandler)
	g.POST("/records/:type", h.CreateRecordHandler)
	g.PUT("/records/:type/:id", h.UpateRecordHandler)
	g.DELETE("/records/:id", h.DeleteRecordHandler)
}

func (o *recordsHandler) GetRecordsHandler(c echo.Context) error {
	userid, err := o.session.GetUserID(c)
	if err != nil {
		o.logger.Errorln(err)
		return c.JSON(http.StatusInternalServerError, httputils.ErrorWithDetails("internal server error", err))
	}

	list, err := o.GetByUser(c.Request().Context(), userid)

	if err != nil {
		o.logger.Errorln(err)
		return c.JSON(http.StatusInternalServerError, httputils.ErrorWithDetails("internal server error", err))
	}

	return c.JSON(http.StatusOK, list)
}

func (o *recordsHandler) UpateRecordHandler(c echo.Context) error {
	ID := c.Param("id")
	dataType := c.Param("type")

	var m interface{}

	switch dataType {
	case record.CREDENTIALS:
		m = new(dto.CredentialsRecord)
	case record.TEXT:
		m = new(dto.TextRecord)
	case record.CREDIT:
		m = new(dto.CreditCardRecord)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "invalid type")
	}

	if err := c.Bind(m); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(m); err != nil {
		return err
	}

	userid, err := o.session.GetUserID(c)
	if err != nil {
		o.logger.Errorln(err)
		return c.JSON(http.StatusInternalServerError, httputils.ErrorWithDetails("internal server error", err))
	}

	err = o.Update(c.Request().Context(), userid, ID, m)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, httputils.ErrorWithDetails("internal server error", err))
	}

	return c.JSON(http.StatusOK, httputils.Ok())
}

func (o *recordsHandler) CreateRecordHandler(c echo.Context) error {
	dataType := c.Param("type")

	var m interface{}
	bindModel := true

	switch dataType {
	case record.CREDENTIALS:
		m = new(dto.CredentialsRecord)
	case record.TEXT:
		m = new(dto.TextRecord)
	case record.BINARY:
		bindModel = false
		sfile, err := c.FormFile("file")
		if err != nil {
			return c.JSON(http.StatusBadRequest, httputils.Error("required file not found"))
		}

		src, err := sfile.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, httputils.ErrorWithDetails("file open error ", err))
		}
		defer src.Close()

		mm := &dto.BinaryRecord{MetaRecord: dto.MetaRecord{Meta: c.FormValue("meta")}, FilePath: c.FormValue("filePath")}

		mm.Data = make([]byte, 0)

		var bf bytes.Buffer
		_, err = io.Copy(&bf, src)

		if err != nil {
			return err
		}
		mm.Data = bf.Bytes()
		m = mm
	case record.CREDIT:
		m = new(dto.CreditCardRecord)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "invalid type")
	}

	if bindModel {
		if err := c.Bind(m); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	if err := c.Validate(m); err != nil {
		return err
	}

	userid, err := o.session.GetUserID(c)
	if err != nil {
		o.logger.Errorln(err)
		return c.JSON(http.StatusInternalServerError, httputils.ErrorWithDetails("user session error", err))
	}

	_, err = o.Create(c.Request().Context(), userid, m, dataType)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, httputils.ErrorWithDetails("internal server error", err))
	}

	return c.JSON(http.StatusOK, httputils.Ok())
}

func (o *recordsHandler) DeleteRecordHandler(c echo.Context) error {
	ID := c.Param("id")

	if len(ID) == 0 {
		return c.JSON(http.StatusBadRequest, errors.New("record not found"))
	}

	userid, err := o.session.GetUserID(c)
	if err != nil {
		o.logger.Errorln(err)
		return c.JSON(http.StatusInternalServerError, httputils.ErrorWithDetails("user session error", err))
	}

	err = o.Delete(c.Request().Context(), userid, ID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, httputils.ErrorWithDetails("internal server error", err))
	}

	return c.JSON(http.StatusOK, httputils.Ok())
}
