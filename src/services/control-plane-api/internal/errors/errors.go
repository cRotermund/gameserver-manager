package errors

import (
	"errors"
	"net/http"
)

var (
	ErrServerNotFound  = errors.New("server not found")
	ErrStateTransition = errors.New("invalid state transition")
	ErrRateLimited     = errors.New("rate limited")
	ErrAwsGeneral      = errors.New("error communicating with AWS")
)

type apiErr struct {
	Message    string
	Code       string
	HTTPStatus int
}

func From(err error) apiErr {
	switch {
	case errors.Is(err, ErrServerNotFound):
		return apiErr{"Resource not found", "NOT_FOUND", http.StatusNotFound}
	case errors.Is(err, ErrStateTransition):
		return apiErr{"Can not transition server state as requested", "STATE_CONFLICT", http.StatusConflict}
	case errors.Is(err, ErrRateLimited):
		return apiErr{"Rate limit exceeded", "RATE_LIMITED", http.StatusTooManyRequests}
	default:
		return apiErr{"Internal server error.", "INTERNAL_ERROR", http.StatusInternalServerError}
	}
}
