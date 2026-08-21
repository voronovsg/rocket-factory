package v1

import (
	"context"

	"github.com/voronovsg/rocket-factory/iam/internal/converter"
	authV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/auth/v1"
)

func (a *api) Whoami(ctx context.Context, req *authV1.WhoamiRequest) (*authV1.WhoamiResponse, error) {
	res, err := a.authService.Whoami(ctx, req.GetSessionUuid())
	if err != nil {
		return nil, err
	}

	return &authV1.WhoamiResponse{
		User:    converter.UserToProto(res.User),
		Session: converter.SessionToProto(res.Session),
	}, nil
}
