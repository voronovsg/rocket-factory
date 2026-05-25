package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/voronovsg/rocket-factory/notification/internal/converter/kafka"
	def "github.com/voronovsg/rocket-factory/notification/internal/service"
	"github.com/voronovsg/rocket-factory/platform/pkg/kafka"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
)

type service struct {
	orderPaidConsumer kafka.Consumer
	orderPaidDecoder  kafkaConverter.OrderPaidDecoder
	telegramService   def.TelegramService
}

func NewService(orderPaidConsumer kafka.Consumer, orderPaidDecoder kafkaConverter.OrderPaidDecoder, telegramService def.TelegramService) *service {
	return &service{
		orderPaidConsumer: orderPaidConsumer,
		orderPaidDecoder:  orderPaidDecoder,
		telegramService:   telegramService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order paid consumer service")

	err := s.orderPaidConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.paid topic error", zap.Error(err))
		return err
	}

	return nil
}
