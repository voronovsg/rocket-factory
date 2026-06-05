package model

import "github.com/pkg/errors"

var (
	ErrUserNotFound               = errors.New("user not found")
	ErrUserUUIDInvalid            = errors.New("user uuid invalid")
	ErrUserLoginOrPasswordInvalid = errors.New("user login or password invalid")
	ErrUserIdentifierInvalid      = errors.New("at least one identifier field must be provided")
	ErrSessionInvalidOrExpired    = errors.New("session invalid or expired")
)
