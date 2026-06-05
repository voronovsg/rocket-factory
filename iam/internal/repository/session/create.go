package session

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/voronovsg/rocket-factory/iam/internal/model"
	"github.com/voronovsg/rocket-factory/iam/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, data model.SessionData, sessionTTL time.Duration) (string, error) {
	sessionUUID := uuid.New().String()
	sessionKey := makeSessionKey(sessionUUID)

	data.Session.UUID = sessionUUID
	data.Session.CreatedAt = time.Now()
	data.Session.ExpiresAt = data.Session.CreatedAt.Add(sessionTTL)

	err := r.rdc.HashSet(
		ctx,
		sessionKey,
		converter.SessionDataToRedisView(ctx, data),
	)
	if err != nil {
		return "", err
	}

	err = r.rdc.Expire(ctx, sessionKey, sessionTTL)
	if err != nil {
		return "", err
	}

	return sessionUUID, nil
}
