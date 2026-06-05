package repository

import (
	"context"
	"time"

	"github.com/voronovsg/rocket-factory/iam/internal/model"
)

type SessionRepository interface {
	Create(ctx context.Context, session model.SessionData, sessionTTL time.Duration) (string, error)
	Get(ctx context.Context, sessionUUID string) (model.SessionData, error)
	AddSessionToUserSet(ctx context.Context, userUUID, sessionUUID string) error
}

type UserRepository interface {
	Create(ctx context.Context, newUser model.UserRegistrationInfo) (string, error)
	GetByIdentifier(ctx context.Context, identifier model.UserIdentifier) (model.User, error)
}
