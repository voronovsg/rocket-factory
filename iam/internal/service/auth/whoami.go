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
		logger.Error(ctx, "ошибка получения данных сессии", zap.Error(err))
		return model.SessionData{}, err
	}

	if sessionData.Session.ExpiresAt.Before(time.Now()) {
		return model.SessionData{}, model.ErrSessionInvalidOrExpired
	}

	return sessionData, nil
}
