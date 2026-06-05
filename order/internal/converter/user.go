package converter

import (
	"time"

	"github.com/voronovsg/rocket-factory/order/internal/model"
	"github.com/voronovsg/rocket-factory/platform/pkg/ptr"
	authV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/auth/v1"
	commonV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/common/v1"
)

func SessionDataToModel(res *authV1.WhoamiResponse) model.SessionData {
	return model.SessionData{
		Session: SessionToModel(res.Session),
		User:    UserToModel(res.User),
	}
}

func UserToModel(user *commonV1.User) model.User {
	if user == nil {
		return model.User{}
	}

	var updatedAt *time.Time
	if user.UpdatedAt != nil {
		updatedAt = ptr.Of(user.UpdatedAt.AsTime())
	}

	return model.User{
		UUID:      user.Uuid,
		Info:      UserInfoToModel(user.Info),
		CreatedAt: user.CreatedAt.AsTime(),
		UpdatedAt: updatedAt,
	}
}

func SessionToModel(session *commonV1.Session) model.Session {
	if session == nil {
		return model.Session{}
	}

	var updatedAt *time.Time
	if session.UpdatedAt != nil {
		updatedAt = ptr.Of(session.UpdatedAt.AsTime())
	}

	return model.Session{
		UUID:      session.Uuid,
		CreatedAt: session.CreatedAt.AsTime(),
		UpdatedAt: updatedAt,
		ExpiresAt: session.ExpiresAt.AsTime(),
	}
}

func UserInfoToModel(info *commonV1.UserInfo) model.UserInfo {
	if info == nil {
		return model.UserInfo{}
	}

	return model.UserInfo{
		Login:               info.Login,
		Email:               info.Email,
		NotificationMethods: NotificationMethodListToModel(info.NotificationMethods),
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
	if method == nil {
		return model.NotificationMethod{}
	}

	return model.NotificationMethod{
		ProviderName: method.ProviderName,
		Target:       method.Target,
	}
}
