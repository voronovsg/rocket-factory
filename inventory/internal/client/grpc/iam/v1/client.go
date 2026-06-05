package v1

import (
	"context"

	def "github.com/voronovsg/rocket-factory/inventory/internal/client/grpc"
	"github.com/voronovsg/rocket-factory/inventory/internal/converter"
	"github.com/voronovsg/rocket-factory/inventory/internal/model"
	authV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/auth/v1"
)

var _ def.IAMClient = (*client)(nil)

type client struct {
	generatedClient authV1.AuthServiceClient
}

// NewClient создает новый IAM клиент
func NewClient(generatedClient authV1.AuthServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}

// Whoami валидирует сессию пользователя
func (c *client) Whoami(ctx context.Context, sessionUUID string) (model.SessionData, error) {
	res, err := c.generatedClient.Whoami(ctx, &authV1.WhoamiRequest{
		SessionUuid: sessionUUID,
	})
	if err != nil {
		return model.SessionData{}, err
	}

	return converter.SessionDataToModel(res), nil
}
