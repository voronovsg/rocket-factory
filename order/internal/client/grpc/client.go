package grpc

import (
	"context"

	"github.com/voronovsg/rocket-factory/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (string, error)
}

type IAMClient interface {
	Whoami(ctx context.Context, sessionUUID string) (model.SessionData, error)
}
