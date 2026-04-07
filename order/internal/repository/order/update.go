package part

import (
	"context"

	"github.com/voronovsg/rocket-factory/order/internal/model"
)

func (r *repository) Update(_ context.Context, orderUUID string, updateOrder model.UpdateOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.orders[orderUUID]
	if !ok {
		return model.ErrOrderNotFound
	}

	if updateOrder.PartUuids != nil {
		order.PartUuids = updateOrder.PartUuids
	}
	if updateOrder.TotalPrice != nil {
		order.TotalPrice = *updateOrder.TotalPrice
	}
	if updateOrder.TransactionUUID != nil {
		order.TransactionUUID = updateOrder.TransactionUUID
	}
	if updateOrder.PaymentMethod != nil {
		order.PaymentMethod = updateOrder.PaymentMethod
	}
	if updateOrder.Status != nil {
		order.Status = *updateOrder.Status
	}

	r.orders[orderUUID] = order

	return nil
}
