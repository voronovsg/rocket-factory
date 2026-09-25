package auth

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/voronovsg/rocket-factory/iam/internal/model"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
)

func (s *service) Whoami(ctx context.Context, sessionUUID string) (model.SessionData, error) {
	sessionData, err := s.sessionRepository.Get(ctx, sessionUUID)
	if err != nil {
		logger.Error(ctx, "Error retrieving session data",
			zap.String("sessionUUID", sessionUUID),
			zap.Error(err))
		return model.SessionData{}, err
	}

	if sessionData.Session.ExpiresAt.Before(time.Now()) {
		logger.Error(ctx, "Session is invalid or expired",
			zap.String("sessionUUID", sessionUUID),
			zap.Error(model.ErrSessionInvalidOrExpired))
		return model.SessionData{}, model.ErrSessionInvalidOrExpired
	}

	logger.Debug(ctx, "Successfully retrieved session data", zap.String("userUUID", sessionData.User.UUID))
	return sessionData, nil
}
