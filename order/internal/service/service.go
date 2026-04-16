package service

import (
	"context"

	"github.com/voronovsg/rocket-factory/order/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, createOrder model.CreateOrder) (model.Order, error)
	GetOrderByUUID(ctx context.Context, orderUUID string) (model.Order, error)
	PayOrder(ctx context.Context, orderUUID, payMethod string) (string, error)
	CancelOrderByUUID(ctx context.Context, orderUUID string) error
}
