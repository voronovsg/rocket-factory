package model

import "github.com/pkg/errors"

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrOrderStatusInvalid = errors.New("invalid order status")
	ErrPartsNotFound      = errors.New("some parts not found")
)
