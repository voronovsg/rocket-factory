package order

import (
	"context"

	"go.uber.org/zap"

	"github.com/voronovsg/rocket-factory/order/internal/model"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
)

func (s *service) GetOrderByUUID(ctx context.Context, orderUUID string) (model.Order, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		logger.Error(ctx, "Failed to get order", zap.String("orderUUID", orderUUID), zap.Error(err))
		return model.Order{}, err
	}

	logger.Debug(ctx, "Get order by UUID", zap.String("orderUUID", order.UUID))
	return order, nil
}
