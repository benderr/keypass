package sender

import (
	"github.com/benderr/keypass/pkg/logger"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	*resty.Client
}

func New(server string, certStr string, logger logger.Logger) *Client {
	client := resty.
		New().
		SetLogger(logger).
		SetBaseURL(server).
		SetRootCertificateFromString(certStr)

	return &Client{
		Client: client,
	}
}
