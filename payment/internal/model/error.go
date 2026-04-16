package model

import "github.com/pkg/errors"

var (
	ErrOrderUUIDInvalid = errors.New("order uuid invalid")
	ErrUserUUIDInvalid  = errors.New("user uuid invalid")
)
