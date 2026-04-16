package v1

import (
	"context"

	"github.com/voronovsg/rocket-factory/order/internal/converter"
	orderV1 "github.com/voronovsg/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrderByUUID(ctx context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	order, err := a.orderService.GetOrderByUUID(ctx, params.OrderUUID.String())
	if err != nil {
		return nil, err
	}

	return &orderV1.GetOrderResponse{
		Order: converter.OrderToProto(order),
	}, nil
}
