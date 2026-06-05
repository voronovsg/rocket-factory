package v1

import (
	"context"

	"github.com/voronovsg/rocket-factory/iam/internal/converter"
	userV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/user/v1"
)

func (a *api) Register(ctx context.Context, req *userV1.RegisterRequest) (*userV1.RegisterResponse, error) {
	userUUID, err := a.userService.Register(ctx, converter.UserRegistrationInfoToModel(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	return &userV1.RegisterResponse{
		UserUuid: userUUID,
	}, nil
}
