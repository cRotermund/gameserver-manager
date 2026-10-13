package apihandling

import "log/slog"

type GsmApiClient struct {
	logger *slog.Logger
}

func NewClient(logger *slog.Logger) *GsmApiClient {

	return &GsmApiClient{
		logger: logger,
	}
}
