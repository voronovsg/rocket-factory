package order_consumer

import (
	"context"
	"math/rand"
	"time"

	"go.uber.org/zap"

	"github.com/voronovsg/rocket-factory/assembly/internal/model"
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

	logger.Info(ctx, "Starting assembly process")
	buildTimeSec := rand.Intn(10) + 1 //nolint:gosec
	sleepTime := time.Duration(buildTimeSec) * time.Second
	timer := time.NewTimer(sleepTime)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
	}
	logger.Info(ctx, "Assembly process completed")

	err = s.orderProducerService.ProduceOrderAssembled(ctx, model.OrderAssembledEvent{
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: int64(buildTimeSec),
	})
	if err != nil {
		logger.Error(ctx, "Failed to republish message", zap.Error(err))
	}

	return nil
}
