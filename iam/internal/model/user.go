package model

import (
	"time"
)

type NotificationMethod struct {
	ProviderName string `json:"provider_name"`
	Target       string `json:"target"`
}

type UserInfo struct {
	Login               string
	Email               string
	NotificationMethods []NotificationMethod
	Password            string
}

type UserRegistrationInfo struct {
	Info UserInfo
}

type User struct {
	UUID      string
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type UserIdentifier struct {
	UUID  *string
	Login *string
	Email *string
}
