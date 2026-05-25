package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/voronovsg/rocket-factory/platform/pkg/kafka/consumer"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
)

func (s *service) OrderHandler(ctx context.Context, msg consumer.Message) error {
	event, err := s.orderAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderAssembled", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("order_uuid", event.OrderUUID),
		zap.Int64("build_time_sec", event.BuildTimeSec),
	)

	err = s.telegramService.SendOrderAssembledNotification(ctx, event)
	if err != nil {
		logger.Error(ctx, "Failed to send assembled notification", zap.Error(err))
		return err
	}

	return nil
}
