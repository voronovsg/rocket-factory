package model

import "github.com/pkg/errors"

var (
	ErrPartNotFound    = errors.New("part not found")
	ErrPartUUIDInvalid = errors.New("part uuid invalid")
)
