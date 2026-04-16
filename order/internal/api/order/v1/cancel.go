package v1

import (
	"context"

	orderV1 "github.com/voronovsg/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrderByUUID(ctx context.Context, params orderV1.CancelOrderByUUIDParams) (orderV1.CancelOrderByUUIDRes, error) {
	if err := a.orderService.CancelOrderByUUID(ctx, params.OrderUUID.String()); err != nil {
		return nil, err
	}

	return &orderV1.CancelOrderByUUIDNoContent{}, nil
}
