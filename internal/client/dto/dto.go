package dto

import (
	"time"

	recordform "github.com/benderr/keypass/pkg/client/component/record_form"
)

type User struct {
	ID    string `json:"id"`
	Login string `json:"login"`
	Token string `json:"token"`
}

// ClientRecord is structure returned from api
type ClientRecord struct {
	ID        string         `json:"id"`
	Meta      string         `json:"meta"`
	Info      map[string]any `json:"info"`
	Version   int            `json:"version"`
	UpdatedAt time.Time      `json:"updated_at"`
	DataType  string         `json:"data_type"`
	UserID    string         `json:"user_id"`
}

// ServerRecord is structure for send to api
type ServerRecord struct {
	ID       string         `json:"id"`
	Meta     string         `json:"meta"`
	Info     map[string]any `json:"info"`
	DataType string         `json:"data_type"`
}

// ConvertToServerRecord convert form values to api payload
func ConvertToServerRecord(rec recordform.RecordValues) ServerRecord {
	v := ServerRecord{
		ID:       rec.ID,
		Meta:     rec.Meta,
		DataType: rec.DataType,
		Info:     make(map[string]any),
	}

	switch rec.DataType {
	case recordform.CREDENTIALS:
		v.Info["login"] = rec.Login
		v.Info["password"] = rec.Password

	case recordform.CREDIT:
		v.Info["number"] = rec.Number
		v.Info["cvv"] = rec.CVV
		v.Info["expire"] = rec.Expire
	case recordform.TEXT:
		v.Info["text"] = rec.Text
	case recordform.BINARY:
		v.Info["filePath"] = rec.FilePath
		v.Info["binary"] = rec.Binary
	}

	return v
}

// ConvertToFormValues convert model from api to form values
func ConvertToFormValues(rec ClientRecord) recordform.RecordValues {
	v := recordform.RecordValues{
		ID:       rec.ID,
		DataType: rec.DataType,
		Meta:     rec.Meta,
	}

	convertToString := func(key string) string {
		data, ok := rec.Info[key]
		if !ok {
			return ""
		}

		if str, ok := data.(string); ok {
			return str
		}
		return ""
	}

	switch rec.DataType {
	case recordform.CREDENTIALS:
		v.Login = convertToString("login")
		v.Password = convertToString("password")
	case recordform.CREDIT:
		v.Number = convertToString("number")
		v.CVV = convertToString("cvv")
		v.Expire = convertToString("expire")
	case recordform.TEXT:
		v.Text = convertToString("text")
	case recordform.BINARY:
		//add file extension
		if data, ok := rec.Info["binary"]; ok {
			if b, ok := data.([]byte); ok {
				v.Binary = b
			}
		}
	}
	return v
}
