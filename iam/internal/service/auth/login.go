package auth

import (
	"context"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/voronovsg/rocket-factory/iam/internal/model"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
)

func (s *service) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.userRepository.GetByIdentifier(ctx, model.UserIdentifier{Login: &login})
	if err != nil {
		logger.Error(ctx, "Failed to get user by login", zap.String("login", login), zap.Error(err))
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Info.Password), []byte(password))
	if err != nil {
		logger.Error(ctx, "Failed to compare password", zap.String("userUUID", user.UUID), zap.Error(err))
		return "", model.ErrUserLoginOrPasswordInvalid
	}

	sessionData := model.SessionData{
		User: user,
	}

	sessionUUID, err := s.sessionRepository.Create(ctx, sessionData, s.sessionConfig.TTL())
	if err != nil {
		logger.Error(ctx, "Failed to create session", zap.String("userUUID", user.UUID), zap.Error(err))
		return "", err
	}

	err = s.sessionRepository.AddSessionToUserSet(ctx, user.UUID, sessionUUID)
	if err != nil {
		logger.Error(ctx, "Failed to add session to user set", zap.String("userUUID", user.UUID), zap.Error(err))
		return "", err
	}

	logger.Debug(ctx, "Successfully logged in", zap.String("user", user.UUID))
	return sessionUUID, nil
}
