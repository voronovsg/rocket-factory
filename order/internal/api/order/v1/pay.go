package v1

import (
	"context"

	"github.com/google/uuid"

	orderV1 "github.com/voronovsg/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	transactionUUID, err := a.orderService.PayOrder(ctx, params.OrderUUID.String(), string(req.PaymentMethod))
	if err != nil {
		return nil, err
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: uuid.MustParse(transactionUUID),
	}, nil
}
