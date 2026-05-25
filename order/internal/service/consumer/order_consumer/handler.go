package order_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/voronovsg/rocket-factory/order/internal/model"
	"github.com/voronovsg/rocket-factory/platform/pkg/kafka/consumer"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
	"github.com/voronovsg/rocket-factory/platform/pkg/ptr"
)

func (s *service) OrderHandler(ctx context.Context, msg consumer.Message) error {
	event, err := s.orderAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Error decoding OrderAssembledEvent event", zap.Error(err))
		return err
	}
	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("order_uuid", event.OrderUUID),
		zap.Int64("build_time_sec", event.BuildTimeSec),
	)
	err = s.orderRepository.Update(ctx, event.OrderUUID, model.UpdateOrder{
		Status: ptr.Of(model.OrderStatusAssembled),
	})
	if err != nil {
		logger.Error(ctx, "Error updating order status", zap.Error(err))
		return err
	}

	return nil
}
