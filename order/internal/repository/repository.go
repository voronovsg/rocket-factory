package repository

import (
	"context"

	"github.com/voronovsg/rocket-factory/order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.CreateOrder) (model.Order, error)
	Get(ctx context.Context, orderUUID string) (model.Order, error)
	Update(ctx context.Context, orderUUID string, updateOrder model.UpdateOrder) error
}
