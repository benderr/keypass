package delivery_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/benderr/keypass/internal/server/domain/record/delivery"
	"github.com/benderr/keypass/internal/server/domain/record/delivery/mocks"
	"github.com/benderr/keypass/internal/server/domain/record/dto"
	mocklogger "github.com/benderr/keypass/pkg/logger/mock_logger"
	"github.com/go-playground/validator"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
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

func newTestServer(mockUsecase delivery.RecordUsecase, mockSession delivery.SessionManager) *httptest.Server {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	mockLogger := mocklogger.New()
	delivery.NewRecordHandlers(e.Group(""), mockUsecase, mockSession, mockLogger)

	return httptest.NewServer(e)
}

func newRequest(baseServer string) *resty.Request {
	return resty.New().SetBaseURL(baseServer).R().SetHeader(echo.HeaderContentType, echo.MIMEApplicationJSON)
}

func TestCreateRecords(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockRecordUsecase(ctrl)
	mockSession := mocks.NewMockSessionManager(ctrl)

	server := newTestServer(mockUsecase, mockSession)
	defer server.Close()

	type want struct {
		code    int
		message string
	}

	tests := []struct {
		name          string
		dataType      string
		userid        string
		meta          string
		requestInfo   map[string]string
		createPayload interface{}
		want          want
	}{
		{
			name:        "Creating success",
			dataType:    "CREDENTIALS",
			userid:      "testuserid",
			meta:        "Test meta",
			requestInfo: map[string]string{"login": "test", "password": "test"},
			createPayload: &dto.CredentialsRecord{
				MetaRecord: dto.MetaRecord{Meta: "Test meta"},
				Info:       dto.CredentialsInfo{Login: "test", Password: "test"},
			},
			want: want{
				code:    http.StatusOK,
				message: `{"message":"ok"}`,
			},
		},
		{
			name:        "Creating invalid contract",
			dataType:    "CREDENTIALS",
			meta:        "Test meta",
			requestInfo: map[string]string{"login": "test", "password": ""},
			want: want{
				code:    http.StatusBadRequest,
				message: `{"message":"Key: 'CredentialsRecord.Info.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
			},
		},
		{
			name:        "Creating success",
			dataType:    "TEXT",
			userid:      "testuserid",
			meta:        "Test meta for text",
			requestInfo: map[string]string{"text": "secret text"},
			createPayload: &dto.TextRecord{
				MetaRecord: dto.MetaRecord{Meta: "Test meta for text"},
				Info:       dto.TextInfo{Text: "secret text"},
			},
			want: want{
				code:    http.StatusOK,
				message: `{"message":"ok"}`,
			},
		},
		{
			name:        "Creating invalid contract",
			dataType:    "TEXT",
			meta:        "Test meta for text",
			requestInfo: map[string]string{"text": ""},
			want: want{
				code:    http.StatusBadRequest,
				message: `{"message":"Key: 'TextRecord.Info.Text' Error:Field validation for 'Text' failed on the 'required' tag"}`,
			},
		},
		{
			name:        "Creating success",
			dataType:    "CREDIT",
			userid:      "testuserid",
			meta:        "Test meta for credit card",
			requestInfo: map[string]string{"number": "123", "cvv": "123", "expire": "12/24"},
			createPayload: &dto.CreditCardRecord{
				MetaRecord: dto.MetaRecord{Meta: "Test meta for credit card"},
				Info:       dto.CreditCardInfo{Number: "123", CVV: "123", Expire: "12/24"},
			},
			want: want{
				code:    http.StatusOK,
				message: `{"message":"ok"}`,
			},
		},
		{
			name:        "Creating invalid contract",
			dataType:    "CREDIT",
			meta:        "Test meta for credit card",
			requestInfo: map[string]string{"cvv": "123", "expire": "12/24"},
			want: want{
				code:    http.StatusBadRequest,
				message: `{"message":"Key: 'CreditCardRecord.Info.Number' Error:Field validation for 'Number' failed on the 'required' tag"}`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name+" for "+test.dataType, func(t *testing.T) {
			if len(test.userid) > 0 {
				mockSession.EXPECT().GetUserID(gomock.Any()).Return(test.userid, nil)
			}

			if test.createPayload != nil {
				mockUsecase.EXPECT().Create(gomock.Any(), test.userid, test.createPayload, test.dataType).Return(true, nil)
			}

			b, _ := json.Marshal(test.requestInfo)

			resp, err := newRequest(server.URL).
				SetBody(fmt.Sprintf(`{"meta": "%v", "info": %v }`, test.meta, string(b))).
				Post("/api/records/" + test.dataType)

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, test.want.code, resp.StatusCode())
			assert.JSONEq(t, test.want.message, string(resp.Body()))
		})
	}
}

func TestCreateFileRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockRecordUsecase(ctrl)
	mockSession := mocks.NewMockSessionManager(ctrl)

	server := newTestServer(mockUsecase, mockSession)
	defer server.Close()

	dataType := "BINARY"

	type want struct {
		code    int
		message string
	}

	tests := []struct {
		name          string
		userid        string
		filePath      string
		meta          string
		createPayload interface{}
		fileBinary    []byte
		want          want
	}{

		{
			name:       "Creating file record success",
			meta:       "Test meta text for file",
			userid:     "testuserid",
			fileBinary: []byte("123 \n 321"),
			filePath:   "c:/user/test/file.txt",
			createPayload: &dto.BinaryRecord{
				MetaRecord: dto.MetaRecord{Meta: "Test meta text for file"},
				Data:       []byte("123 \n 321"),
				FilePath:   "c:/user/test/file.txt",
			},
			want: want{
				code:    http.StatusOK,
				message: `{"message":"ok"}`,
			},
		},
		{
			name:     "Creating file record error",
			meta:     "Test meta text for file",
			filePath: "c:/user/test/file.txt",
			want: want{
				code:    http.StatusBadRequest,
				message: `{"message":"required file not found"}`,
			},
		},
		{
			name:       "Creating file record error, no meta",
			fileBinary: []byte("123 \n 321"),
			filePath:   "c:/user/test/file.txt",
			want: want{
				code:    http.StatusBadRequest,
				message: `{"message":"Key: 'BinaryRecord.MetaRecord.Meta' Error:Field validation for 'Meta' failed on the 'required' tag"}`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if len(test.userid) > 0 {
				mockSession.EXPECT().GetUserID(gomock.Any()).Return(test.userid, nil)
			}

			if test.createPayload != nil {
				mockUsecase.EXPECT().Create(gomock.Any(), test.userid, test.createPayload, dataType).Return(true, nil)
			}

			r := newRequest(server.URL)

			if len(test.fileBinary) > 0 {
				r.SetFileReader("file", test.filePath, bytes.NewReader(test.fileBinary))
			}

			r.SetFormData(map[string]string{
				"meta":     test.meta,
				"filePath": test.filePath,
			})

			resp, err := r.Post("/api/records/" + dataType)

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, test.want.code, resp.StatusCode())
			assert.JSONEq(t, test.want.message, string(resp.Body()))
		})
	}
}

func TestUpdateRecords(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockRecordUsecase(ctrl)
	mockSession := mocks.NewMockSessionManager(ctrl)

	server := newTestServer(mockUsecase, mockSession)
	defer server.Close()

	type want struct {
		code    int
		message string
	}

	tests := []struct {
		ID            string
		name          string
		dataType      string
		userid        string
		meta          string
		requestInfo   map[string]string
		createPayload interface{}
		want          want
	}{
		{
			ID:          "123",
			name:        "Updating success",
			dataType:    "CREDENTIALS",
			userid:      "testuserid",
			meta:        "Test meta",
			requestInfo: map[string]string{"login": "test", "password": "test"},
			createPayload: &dto.CredentialsRecord{
				MetaRecord: dto.MetaRecord{Meta: "Test meta"},
				Info:       dto.CredentialsInfo{Login: "test", Password: "test"},
			},
			want: want{
				code:    http.StatusOK,
				message: `{"message":"ok"}`,
			},
		},
		{
			ID:          "123",
			name:        "Updating invalid contract",
			dataType:    "CREDENTIALS",
			meta:        "Test meta",
			requestInfo: map[string]string{"login": "test", "password": ""},
			want: want{
				code:    http.StatusBadRequest,
				message: `{"message":"Key: 'CredentialsRecord.Info.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
			},
		},
		{
			ID:          "123",
			name:        "Updating success",
			dataType:    "TEXT",
			userid:      "testuserid",
			meta:        "Test meta for text",
			requestInfo: map[string]string{"text": "secret text"},
			createPayload: &dto.TextRecord{
				MetaRecord: dto.MetaRecord{Meta: "Test meta for text"},
				Info:       dto.TextInfo{Text: "secret text"},
			},
			want: want{
				code:    http.StatusOK,
				message: `{"message":"ok"}`,
			},
		},
		{
			ID:          "123",
			name:        "Updating invalid contract",
			dataType:    "TEXT",
			meta:        "Test meta for text",
			requestInfo: map[string]string{"text": ""},
			want: want{
				code:    http.StatusBadRequest,
				message: `{"message":"Key: 'TextRecord.Info.Text' Error:Field validation for 'Text' failed on the 'required' tag"}`,
			},
		},
		{
			ID:          "123",
			name:        "Updating success",
			dataType:    "CREDIT",
			userid:      "testuserid",
			meta:        "Test meta for credit card",
			requestInfo: map[string]string{"number": "123", "cvv": "123", "expire": "12/24"},
			createPayload: &dto.CreditCardRecord{
				MetaRecord: dto.MetaRecord{Meta: "Test meta for credit card"},
				Info:       dto.CreditCardInfo{Number: "123", CVV: "123", Expire: "12/24"},
			},
			want: want{
				code:    http.StatusOK,
				message: `{"message":"ok"}`,
			},
		},
		{
			ID:          "123",
			name:        "Updating invalid contract",
			dataType:    "CREDIT",
			meta:        "Test meta for credit card",
			requestInfo: map[string]string{"cvv": "123", "expire": "12/24"},
			want: want{
				code:    http.StatusBadRequest,
				message: `{"message":"Key: 'CreditCardRecord.Info.Number' Error:Field validation for 'Number' failed on the 'required' tag"}`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name+" for "+test.dataType, func(t *testing.T) {
			if len(test.userid) > 0 {
				mockSession.EXPECT().GetUserID(gomock.Any()).Return(test.userid, nil)
			}

			if test.createPayload != nil {
				mockUsecase.EXPECT().Update(gomock.Any(), test.userid, test.ID, test.createPayload).Return(nil)
			}

			b, _ := json.Marshal(test.requestInfo)

			resp, err := newRequest(server.URL).
				SetBody(fmt.Sprintf(`{"meta": "%v", "info": %v }`, test.meta, string(b))).
				Put("/api/records/" + test.dataType + "/" + test.ID)

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, test.want.code, resp.StatusCode())
			assert.JSONEq(t, test.want.message, string(resp.Body()))
		})
	}
}

func TestGetRecordsRecords(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockRecordUsecase(ctrl)
	mockSession := mocks.NewMockSessionManager(ctrl)

	server := newTestServer(mockUsecase, mockSession)
	defer server.Close()

	records := make([]dto.ReadRecord, 0)
	records = append(records, dto.ReadRecord{
		ID:        "1",
		Meta:      "Meta",
		Info:      map[string]any{},
		Version:   1,
		UpdatedAt: time.Date(2010, 1, 1, 1, 1, 1, 1, time.UTC),
		DataType:  "CREDIT",
	})

	responseJson, _ := json.Marshal(records)

	type want struct {
		code    int
		message string
	}

	tests := []struct {
		name         string
		userid       string
		userError    error
		recordsError error
		records      interface{}
		want         want
	}{
		{
			name:         "Get records success",
			userid:       "testuserid",
			records:      records,
			userError:    nil,
			recordsError: nil,
			want: want{
				code:    http.StatusOK,
				message: string(responseJson),
			},
		},
		{
			name:      "Get records failed by user",
			userError: errors.New("user not found"),
			want: want{
				code:    http.StatusInternalServerError,
				message: `{"details":"user not found", "message":"internal server error"}`,
			},
		},
		{
			name:         "Get records failed by record-usecase",
			userid:       "testuserid",
			recordsError: errors.New("fetch records error"),
			want: want{
				code:    http.StatusInternalServerError,
				message: `{"details":"fetch records error", "message":"internal server error"}`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if len(test.userid) > 0 || test.userError != nil {
				mockSession.EXPECT().GetUserID(gomock.Any()).Return(test.userid, test.userError)
			}

			if test.records != nil || test.recordsError != nil {
				mockUsecase.EXPECT().GetByUser(gomock.Any(), test.userid).Return(test.records, test.recordsError)
			}

			resp, err := newRequest(server.URL).Get("/api/records")

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, test.want.code, resp.StatusCode())
			if len(test.want.message) > 0 {
				assert.JSONEq(t, test.want.message, string(resp.Body()))
			}
		})
	}
}

func TestDeleteRecords(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockRecordUsecase(ctrl)
	mockSession := mocks.NewMockSessionManager(ctrl)

	server := newTestServer(mockUsecase, mockSession)
	defer server.Close()

	userid := "test_user"
	ID := "record_id"

	mockSession.EXPECT().GetUserID(gomock.Any()).Return(userid, nil)

	mockUsecase.EXPECT().Delete(gomock.Any(), userid, ID).Return(nil)

	resp, err := newRequest(server.URL).Delete("/api/records/" + ID)

	assert.NoError(t, err, "error making HTTP request")

	assert.Equal(t, http.StatusOK, resp.StatusCode())
}
