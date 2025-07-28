package hsr

import "errors"

var (
	ErrInvalidUIDFormat   = errors.New("invalid UID format")
	ErrPlayerNotFound     = errors.New("player not found")
	ErrServerMaintenance  = errors.New("server maintenance")
	ErrServerError        = errors.New("server error")
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrRateLimited        = errors.New("rate limited")
)
