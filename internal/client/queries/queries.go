package queries

import (
	"bytes"
	"errors"

	"github.com/benderr/keypass/internal/client/dto"
	"github.com/benderr/keypass/pkg/sender"
)

type queryClient struct {
	client *sender.Client
}

func New(c *sender.Client) *queryClient {
	return &queryClient{
		client: c,
	}
}

type AuthModel struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ResponseError struct {
	Message string `json:"message"`
}

func (r *ResponseError) Error() string {
	return r.Message
}

func (q *queryClient) Login(login string, pass string) (dto.User, error) {
	u := &dto.User{}

	resp, err := q.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(AuthModel{Login: login, Password: pass}).
		SetResult(u).
		SetError(ResponseError{}).
		Post("/api/user/login")

	if err, ok := resp.Error().(*ResponseError); ok {
		return *u, err
	}

	if err != nil {
		return *u, err
	}

	return *u, nil
}

func (q *queryClient) Register(login string, pass string) (dto.User, error) {
	u := &dto.User{}
	resp, err := q.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(AuthModel{login, pass}).
		SetError(ResponseError{}).
		SetResult(u).
		Post("/api/user/register")

	if err != nil {
		return *u, err
	}

	if err, ok := resp.Error().(*ResponseError); ok {
		return *u, err
	}

	return *u, nil
}

func (q *queryClient) GetRecords(token string) ([]dto.ClientRecord, error) {
	u := []dto.ClientRecord{}
	resp, err := q.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Encoding", "gzip").
		SetError(ResponseError{}).
		SetResult(&u).
		SetAuthToken(token).
		Get("/api/records")

	if err != nil {
		return u, errors.New("offline")
	}

	if err, ok := resp.Error().(*ResponseError); ok {
		return u, err
	}

	return u, nil
}

func (q *queryClient) UpdateRecord(token string, record dto.ServerRecord) error {
	resp, err := q.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(record).
		SetError(ResponseError{}).
		SetAuthToken(token).
		Put("/api/records/" + record.DataType + "/" + record.ID)

	if err != nil {
		return err
	}

	if err, ok := resp.Error().(*ResponseError); ok {
		return err
	}

	return nil
}

func (q *queryClient) AddRecordFile(token string, record dto.ServerRecord) error {
	filePath, ok := record.Info["filePath"].(string)
	if !ok {
		return errors.New("invalid filePath")
	}
	fileBinary, ok := record.Info["binary"].([]byte)

	if !ok {
		return errors.New("invalid file content")
	}

	resp, err := q.client.R().
		SetFileReader("file", filePath, bytes.NewReader(fileBinary)).
		SetFormData(map[string]string{
			"meta":     record.Meta,
			"filePath": filePath,
		}).
		SetError(ResponseError{}).
		SetAuthToken(token).
		Post("/api/records/" + record.DataType)

	if err != nil {
		return err
	}

	if err, ok := resp.Error().(*ResponseError); ok {
		return err
	}

	return nil
}

func (q *queryClient) AddRecord(token string, record dto.ServerRecord) error {
	resp, err := q.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(record).
		SetError(ResponseError{}).
		SetAuthToken(token).
		Post("/api/records/" + record.DataType)

	if err != nil {
		return err
	}

	if err, ok := resp.Error().(*ResponseError); ok {
		return err
	}

	return nil
}

func (q *queryClient) DeleteRecord(token string, ID string) error {
	resp, err := q.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Encoding", "gzip").
		SetError(ResponseError{}).
		SetAuthToken(token).
		Delete("/api/records/" + ID)

	if err != nil {
		return err
	}

	if err, ok := resp.Error().(*ResponseError); ok {
		return err
	}

	return nil
}
