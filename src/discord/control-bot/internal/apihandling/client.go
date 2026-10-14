package apihandling

import (
	"log/slog"

	"resty.dev/v3"
)

type GsmApiClient struct {
	logger *slog.Logger
	host   string
	port   int
	rc     *resty.Client
}

func NewClient(logger *slog.Logger, host string, port int, rc *resty.Client) *GsmApiClient {

	return &GsmApiClient{
		logger: logger,
		host:   host,
		port:   port,
		rc:     rc,
	}
}

func (c *GsmApiClient) Healthcheck() {

}

func (c *GsmApiClient) Close() {
	c.rc.Close()
}
