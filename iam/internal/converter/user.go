package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/voronovsg/rocket-factory/iam/internal/model"
	commonV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/common/v1"
	userV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/user/v1"
)

func UserInfoToModel(info *commonV1.UserInfo, password string) model.UserInfo {
	return model.UserInfo{
		Login:               info.Login,
		Email:               info.Email,
		NotificationMethods: NotificationMethodListToModel(info.NotificationMethods),
		Password:            password,
	}
}

func UserRegistrationInfoToModel(info *userV1.UserRegistrationInfo) model.UserRegistrationInfo {
	return model.UserRegistrationInfo{
		Info: UserInfoToModel(info.GetInfo(), info.Password),
	}
}

func UserToProto(user model.User) *commonV1.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt != nil {
		updatedAt = timestamppb.New(*user.UpdatedAt)
	}

	return &commonV1.User{
		Uuid: user.UUID,
		Info: &commonV1.UserInfo{
			Login:               user.Info.Login,
			Email:               user.Info.Email,
			NotificationMethods: NotificationMethodListToProto(user.Info.NotificationMethods),
		},
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func NotificationMethodListToProto(methods []model.NotificationMethod) []*commonV1.NotificationMethod {
	res := make([]*commonV1.NotificationMethod, 0, len(methods))
	for _, method := range methods {
		res = append(res, NotificationMethodToProto(method))
	}

	return res
}

func NotificationMethodToProto(method model.NotificationMethod) *commonV1.NotificationMethod {
	return &commonV1.NotificationMethod{
		ProviderName: method.ProviderName,
		Target:       method.Target,
	}
}

func NotificationMethodListToModel(methods []*commonV1.NotificationMethod) []model.NotificationMethod {
	res := make([]model.NotificationMethod, 0, len(methods))
	for _, method := range methods {
		res = append(res, NotificationMethodToModel(method))
	}

	return res
}

func NotificationMethodToModel(method *commonV1.NotificationMethod) model.NotificationMethod {
	return model.NotificationMethod{
		ProviderName: method.ProviderName,
		Target:       method.Target,
	}
}
