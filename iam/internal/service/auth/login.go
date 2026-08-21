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
		logger.Error(ctx, "failed to get user by login", zap.Error(err))
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Info.Password), []byte(password))
	if err != nil {
		logger.Error(ctx, "failed to compare password", zap.Error(err))
		return "", model.ErrUserLoginOrPasswordInvalid
	}

	sessionData := model.SessionData{
		User: user,
	}

	sessionUUID, err := s.sessionRepository.Create(ctx, sessionData, s.sessionConfig.TTL())
	if err != nil {
		logger.Error(ctx, "failed to create session", zap.Error(err))
		return "", err
	}

	err = s.sessionRepository.AddSessionToUserSet(ctx, user.UUID, sessionUUID)
	if err != nil {
		logger.Error(ctx, "failed to add session to user set", zap.Error(err))
		return "", err
	}

	return sessionUUID, nil
}
