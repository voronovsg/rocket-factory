package v1

import (
	"context"

	"github.com/voronovsg/rocket-factory/iam/internal/converter"
	userV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/user/v1"
)

func (a *api) GetUser(ctx context.Context, req *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	user, err := a.userService.Get(ctx, req.GetUserUuid())
	if err != nil {
		return nil, err
	}

	return &userV1.GetUserResponse{
		User: converter.UserToProto(user),
	}, nil
}
