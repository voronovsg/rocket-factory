package order

import (
	"context"
	"log"

	"github.com/samber/lo"

	"github.com/voronovsg/rocket-factory/order/internal/model"
)

func (s *service) CancelOrderByUUID(ctx context.Context, orderUUID string) error {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return err
	}

	if order.Status != model.OrderStatusPendingPayment {
		return model.ErrOrderStatusInvalid
	}

	err = s.orderRepository.Update(ctx, orderUUID, model.UpdateOrder{
		Status: lo.ToPtr(model.OrderStatusCancelled),
	})
	if err != nil {
		return err
	}

	log.Printf(`[Order Canceled]
	Order UUID: %s
	User UUID: %s
	Status: %s`, orderUUID, order.UserUUID, order.Status)

	return nil
}
