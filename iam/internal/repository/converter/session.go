package converter

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"

	"github.com/voronovsg/rocket-factory/iam/internal/model"
	repoModel "github.com/voronovsg/rocket-factory/iam/internal/repository/model"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
)

func SessionDataToRedisView(ctx context.Context, data model.SessionData) repoModel.SessionDataRedisView {
	var sessionUpdatedAtNs *int64
	if data.Session.UpdatedAt != nil {
		tmp := data.Session.UpdatedAt.UnixNano()
		sessionUpdatedAtNs = &tmp
	}

	var userUpdatedAtNs *int64
	if data.User.UpdatedAt != nil {
		tmp := data.User.UpdatedAt.UnixNano()
		userUpdatedAtNs = &tmp
	}

	var userNotificationMethods string
	if data.User.Info.NotificationMethods != nil {
		raw, err := json.Marshal(data.User.Info.NotificationMethods)
		if err != nil {
			logger.Error(ctx, "ошибка при парсинге notificationMethods", zap.Error(err))
		}

		userNotificationMethods = string(raw)
	}

	return repoModel.SessionDataRedisView{
		SessionUUID:             data.Session.UUID,
		SessionCreatedAtNs:      data.Session.CreatedAt.UnixNano(),
		SessionUpdatedAtNs:      sessionUpdatedAtNs,
		SessionExpiresAtNs:      data.Session.ExpiresAt.UnixNano(),
		UserUUID:                data.User.UUID,
		UserLogin:               data.User.Info.Login,
		UserEmail:               data.User.Info.Email,
		UserNotificationMethods: userNotificationMethods,
		UserCreatedAtAtNs:       data.User.CreatedAt.UnixNano(),
		UserUpdatedAtAtNs:       userUpdatedAtNs,
	}
}

func SessionDataToModel(ctx context.Context, data repoModel.SessionDataRedisView) model.SessionData {
	var sessionUpdatedAt *time.Time
	if data.SessionUpdatedAtNs != nil {
		tmp := time.Unix(0, *data.SessionUpdatedAtNs)
		sessionUpdatedAt = &tmp
	}

	var userUpdatedAt *time.Time
	if data.UserUpdatedAtAtNs != nil {
		tmp := time.Unix(0, *data.UserUpdatedAtAtNs)
		userUpdatedAt = &tmp
	}

	var notificationMethods []model.NotificationMethod
	if data.UserNotificationMethods != "" {
		err := json.Unmarshal([]byte(data.UserNotificationMethods), &notificationMethods)
		if err != nil {
			logger.Error(ctx, "ошибка при парсинге notificationMethods", zap.Error(err))
		}
	}

	return model.SessionData{
		Session: model.Session{
			UUID:      data.SessionUUID,
			CreatedAt: time.Unix(0, data.SessionCreatedAtNs),
			UpdatedAt: sessionUpdatedAt,
			ExpiresAt: time.Unix(0, data.SessionExpiresAtNs),
		},
		User: model.User{
			UUID: data.UserUUID,
			Info: model.UserInfo{
				Login:               data.UserLogin,
				Email:               data.UserEmail,
				NotificationMethods: notificationMethods,
			},
			CreatedAt: time.Unix(0, data.UserCreatedAtAtNs),
			UpdatedAt: userUpdatedAt,
		},
	}
}
