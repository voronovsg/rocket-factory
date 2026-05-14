package order

import (
	"context"
	"errors"
	"log"

	"github.com/voronovsg/rocket-factory/order/internal/model"
	"github.com/voronovsg/rocket-factory/platform/pkg/ptr"
)

func (s *service) PayOrder(ctx context.Context, orderUUID, paymentMethod string) (string, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return "", err
	}

	if order.Status != model.OrderStatusPendingPayment {
		return "", model.ErrOrderStatusInvalid
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, orderUUID, order.UserUUID, paymentMethod)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return "", errors.New("payment service timeout")
		}

		return "", err
	}

	err = s.orderRepository.Update(ctx, orderUUID, model.UpdateOrder{
		TransactionUUID: &transactionUUID,
		PaymentMethod:   &paymentMethod,
		Status:          ptr.Of(model.OrderStatusPaid),
	})
	if err != nil {
		return "", err
	}

	log.Printf(`[Order Paid]
	Order UUID: %s
	User UUID: %s
	Transaction UUID: %s
	Payment Method: %s
	Status: %s`, order.UUID, order.UserUUID, transactionUUID, paymentMethod, model.OrderStatusPaid)

	return transactionUUID, nil
}
