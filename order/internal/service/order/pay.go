package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/voronovsg/rocket-factory/order/internal/model"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
	"github.com/voronovsg/rocket-factory/platform/pkg/ptr"
)

func (s *service) PayOrder(ctx context.Context, orderUUID, paymentMethod string) (string, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		logger.Error(ctx, "failed to get order", zap.Error(err), zap.String("orderUUID", orderUUID))
		return "", err
	}

	if order.Status != model.OrderStatusPendingPayment {
		logger.Error(
			ctx,
			"order status is not pending payment",
			zap.String("status", order.Status),
			zap.String("orderUUID", orderUUID))
		return "", model.ErrOrderStatusInvalid
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, orderUUID, order.UserUUID, paymentMethod)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return "", errors.New("payment service timeout")
		}
		logger.Error(ctx, "failed to pay order", zap.Error(err), zap.String("orderUUID", orderUUID))
		return "", err
	}

	err = s.orderRepository.Update(ctx, orderUUID, model.UpdateOrder{
		TransactionUUID: &transactionUUID,
		PaymentMethod:   &paymentMethod,
		Status:          ptr.Of(model.OrderStatusPaid),
	})
	if err != nil {
		logger.Error(ctx, "failed to update order", zap.Error(err), zap.String("orderUUID", orderUUID))
		return "", err
	}

	err = s.orderProducerService.ProduceOrderPaid(ctx, model.OrderPaidEvent{
		OrderUUID:       order.UUID,
		UserUUID:        order.UserUUID,
		PaymentMethod:   paymentMethod,
		TransactionUUID: transactionUUID,
		PaidAt:          time.Now(),
	})
	if err != nil {
		logger.Error(ctx, "failed to send OrderPaidEvent", zap.Error(err), zap.String("orderUUID", orderUUID))
		return "", err
	}

	logger.Info(ctx, fmt.Sprintf(`[Order Paid]
	Order UUID: %s
	User UUID: %s
	Transaction UUID: %s
	Payment Method: %s
	Status: %s`, order.UUID, order.UserUUID, transactionUUID, paymentMethod, model.OrderStatusPaid))

	return transactionUUID, nil
}
