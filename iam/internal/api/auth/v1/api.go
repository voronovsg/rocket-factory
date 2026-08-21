package v1

import (
	"github.com/voronovsg/rocket-factory/iam/internal/service"
	authV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/auth/v1"
)

type api struct {
	authV1.UnimplementedAuthServiceServer

	authService service.AuthService
}

func NewAPI(authService service.AuthService) *api {
	return &api{
		authService: authService,
	}
}
