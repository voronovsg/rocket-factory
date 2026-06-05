package converter

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"

	"github.com/voronovsg/rocket-factory/iam/internal/model"
	repoModel "github.com/voronovsg/rocket-factory/iam/internal/repository/model"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
)

func UserRegistrationToRepoModel(ctx context.Context, in model.UserRegistrationInfo) repoModel.UserRegistrationInfo {
	notificationMethods, err := json.Marshal(in.Info.NotificationMethods)
	if err != nil {
		logger.Error(ctx, "failed to marshal notification methods", zap.Error(err))
	}

	return repoModel.UserRegistrationInfo{
		Login:               in.Info.Login,
		Email:               in.Info.Email,
		Password:            in.Info.Password,
		NotificationMethods: notificationMethods,
	}
}

func UserToModel(ctx context.Context, user repoModel.User) model.User {
	var notificationMethods []model.NotificationMethod
	err := json.Unmarshal(user.NotificationMethods, &notificationMethods)
	if err != nil {
		logger.Error(ctx, "failed to unmarshal notification methods", zap.Error(err))
	}

	return model.User{
		UUID: user.UUID,
		Info: model.UserInfo{
			Login:               user.Login,
			Email:               user.Email,
			NotificationMethods: notificationMethods,
			Password:            user.Password,
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
