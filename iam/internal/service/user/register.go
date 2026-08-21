package user

import (
	"context"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/voronovsg/rocket-factory/iam/internal/model"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
)

func (s *service) Register(ctx context.Context, user model.UserRegistrationInfo) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Info.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(ctx, "failed to hash password", zap.Error(err))
		return "", err
	}

	user.Info.Password = string(hashedPassword)

	userUUID, err := s.userRepository.Create(ctx, user)
	if err != nil {
		logger.Error(ctx, "failed to create user", zap.Error(err))
		return "", err
	}

	return userUUID, nil
}
