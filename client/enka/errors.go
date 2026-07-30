package enka

import (
	"errors"

	"github.com/kirinyoku/enkanetwork-go/internal/core"
)

var (
	ErrInvalidUsername           = errors.New("username cannot be empty")
	ErrUserNotFound              = errors.New("user not found")
	ErrHoyoAccountNotFound       = errors.New("hoyo account not found")
	ErrHoyoAccountBuildsNotFound = errors.New("no builds found for hoyo account")
	ErrInvalidHoyoHash           = errors.New("hoyo_hash cannot be empty")
	ErrServerMaintenance         = core.ErrServerMaintenance
	ErrServerError               = core.ErrServerError
	ErrServiceUnavailable        = core.ErrServiceUnavailable
	ErrRateLimited               = core.ErrRateLimited
)
