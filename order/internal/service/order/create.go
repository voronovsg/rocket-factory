package order

import (
	"context"
	"time"

	"github.com/go-faster/errors"
	"go.uber.org/zap"

	"github.com/voronovsg/rocket-factory/order/internal/metrics"
	"github.com/voronovsg/rocket-factory/order/internal/model"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
)

func (s *service) CreateOrder(ctx context.Context, createOrder model.CreateOrder) (model.Order, error) {
	parts, err := s.inventoryClient.ListParts(ctx, model.PartsFilter{
		Uuids: createOrder.PartUuids,
	})
	if err != nil {
		logger.Error(ctx, "Inventory service error", zap.Error(err))
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return model.Order{}, errors.New("inventory service timeout")
		}
		return model.Order{}, errors.Errorf("inventory service error: %s", err.Error())
	}

	if len(parts) != len(createOrder.PartUuids) {
		logger.Error(ctx, "Some parts not found",
			zap.Strings("partUUIDs", createOrder.PartUuids),
			zap.Int("expected nums", len(createOrder.PartUuids)),
			zap.Int("actual nums", len(parts)),
			zap.Error(model.ErrPartsNotFound))
		return model.Order{}, model.ErrPartsNotFound
	}

	var totalPrice float64
	for _, part := range parts {
		totalPrice += part.Price
	}

	createOrder.TotalPrice = totalPrice
	createOrder.Status = model.OrderStatusPendingPayment

	order, err := s.orderRepository.Create(ctx, createOrder)
	if err != nil {
		logger.Error(ctx, "Failed to create order", zap.Error(err))
		return model.Order{}, err
	}

	metrics.CountOrderCreated(ctx, model.OrderStatusPendingPayment)

	logger.Debug(ctx, "Order created",
		zap.String("orderUUID", order.UUID),
		zap.String("userUUID", order.UserUUID),
		zap.Strings("partsUUID", order.PartUuids),
		zap.Float64("totalPrice", order.TotalPrice),
		zap.String("status", order.Status),
		zap.String("createdAt", order.CreatedAt.Format(time.RFC3339)))

	return order, nil
}
