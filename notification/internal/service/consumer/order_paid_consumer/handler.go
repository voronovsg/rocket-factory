package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/voronovsg/rocket-factory/platform/pkg/kafka/consumer"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
)

func (s *service) OrderHandler(ctx context.Context, msg consumer.Message) error {
	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.String("transaction_uuid", event.TransactionUUID),
		zap.String("paid_at", event.PaidAt.String()),
	)

	err = s.telegramService.SendOrderPaidNotification(ctx, event)
	if err != nil {
		logger.Error(ctx, "Failed to send paid notification", zap.Error(err))
		return err
	}

	return nil
}
