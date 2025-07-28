package enka

import "errors"

var (
	ErrInvalidUsername           = errors.New("username cannot be empty")
	ErrUserNotFound              = errors.New("user not found")
	ErrHoyoAccountNotFound       = errors.New("hoyo account not found")
	ErrHoyoAccountBuildsNotFound = errors.New("no builds found for hoyo account")
	ErrInvalidHoyoHash           = errors.New("hoyo_hash cannot be empty")
)
