package service

import (
	"context"

	"github.com/voronovsg/rocket-factory/iam/internal/model"
)

type AuthService interface {
	Login(ctx context.Context, login, password string) (string, error)
	Whoami(ctx context.Context, sessionUUID string) (model.SessionData, error)
}

type UserService interface {
	Register(ctx context.Context, user model.UserRegistrationInfo) (string, error)
	Get(ctx context.Context, userUUID string) (model.User, error)
}
