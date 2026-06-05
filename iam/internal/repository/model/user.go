package model

import "time"

type User struct {
	UUID                string     `db:"uuid"`
	Login               string     `db:"login"`
	Email               string     `db:"email"`
	Password            string     `db:"password"`
	NotificationMethods []byte     `db:"notification_methods"`
	CreatedAt           time.Time  `db:"created_at"`
	UpdatedAt           *time.Time `db:"updated_at"`
}

type UserRegistrationInfo struct {
	Login               string
	Email               string
	Password            string
	NotificationMethods []byte
}
