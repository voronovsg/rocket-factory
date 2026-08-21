package model

import (
	"time"
)

type Session struct {
	UUID      string
	CreatedAt time.Time
	UpdatedAt *time.Time
	ExpiresAt time.Time
}

type SessionData struct {
	Session Session
	User    User
}
