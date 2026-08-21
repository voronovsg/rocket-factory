package model

import "time"

type NotificationMethod struct {
	ProviderName string
	Target       string
}

type UserInfo struct {
	Login               string
	Email               string
	NotificationMethods []NotificationMethod
}

type User struct {
	UUID      string
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt *time.Time
}
