package model

type SessionDataRedisView struct {
	SessionUUID        string `redis:"session_uuid"`
	SessionCreatedAtNs int64  `redis:"session_created_at"`
	SessionUpdatedAtNs *int64 `redis:"session_updated_at"`
	SessionExpiresAtNs int64  `redis:"session_expires_at"`

	UserUUID                string `redis:"user_uuid"`
	UserLogin               string `redis:"user_login"`
	UserEmail               string `redis:"user_email"`
	UserNotificationMethods string `redis:"user_notification_methods"`
	UserCreatedAtAtNs       int64  `redis:"user_created_at"`
	UserUpdatedAtAtNs       *int64 `redis:"user_updated_at"`
}
