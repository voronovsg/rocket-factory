package grpc

import (
	"context"

	"github.com/voronovsg/rocket-factory/inventory/internal/model"
)

// IAMClient интерфейс для взаимодействия с IAM сервисом
type IAMClient interface {
	Whoami(ctx context.Context, sessionUUID string) (model.SessionData, error)
}
