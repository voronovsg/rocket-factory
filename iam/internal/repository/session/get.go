package session

import (
	"context"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	"github.com/voronovsg/rocket-factory/iam/internal/model"
	repoConverter "github.com/voronovsg/rocket-factory/iam/internal/repository/converter"
	repoModel "github.com/voronovsg/rocket-factory/iam/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, sessionUUID string) (model.SessionData, error) {
	sessionKey := makeSessionKey(sessionUUID)

	values, err := r.rdc.HGetAll(ctx, sessionKey)
	if err != nil {
		if errors.Is(err, redigo.ErrNil) {
			return model.SessionData{}, model.ErrSessionInvalidOrExpired
		}

		return model.SessionData{}, err
	}
	if len(values) == 0 {
		return model.SessionData{}, model.ErrSessionInvalidOrExpired
	}

	var sessionDataRedisView repoModel.SessionDataRedisView
	err = redigo.ScanStruct(values, &sessionDataRedisView)
	if err != nil {
		return model.SessionData{}, err
	}

	return repoConverter.SessionDataToModel(ctx, sessionDataRedisView), nil
}
