package order

import (
	"context"

	"go.uber.org/zap"

	"github.com/voronovsg/rocket-factory/order/internal/metrics"
	"github.com/voronovsg/rocket-factory/order/internal/model"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
	"github.com/voronovsg/rocket-factory/platform/pkg/ptr"
)

func (s *service) CancelOrderByUUID(ctx context.Context, orderUUID string) error {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		logger.Error(ctx, "Failed to get order", zap.String("orderUUID", orderUUID), zap.Error(err))
		return err
	}

	if order.Status != model.OrderStatusPendingPayment {
		logger.Error(ctx, "Order status is not pending payment",
			zap.String("status", order.Status),
			zap.String("orderUUID", orderUUID))
		return model.ErrOrderStatusInvalid
	}

	err = s.orderRepository.Update(ctx, orderUUID, model.UpdateOrder{
		Status: ptr.Of(model.OrderStatusCancelled),
	})
	if err != nil {
		logger.Error(ctx, "Failed to cancel order", zap.String("orderUUID", orderUUID), zap.Error(err))
		return err
	}

	metrics.CountOrderCreated(ctx, model.OrderStatusCancelled)

	logger.Debug(ctx, "Cancelled order",
		zap.String("orderUUID", orderUUID),
		zap.String("orderUserUUID", order.UserUUID),
		zap.String("status", order.Status))

	return nil
}
