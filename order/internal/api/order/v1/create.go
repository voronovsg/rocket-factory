package v1

import (
	"context"

	"github.com/voronovsg/rocket-factory/order/internal/converter"
	orderV1 "github.com/voronovsg/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest, _ orderV1.CreateOrderParams) (orderV1.CreateOrderRes, error) {
	order, err := a.orderService.CreateOrder(ctx, converter.CreateOrderToModel(req))
	if err != nil {
		return nil, err
	}

	return &orderV1.CreateOrderResponse{
		Order: converter.OrderToProto(order),
	}, nil
}
