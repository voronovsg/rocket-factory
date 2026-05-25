package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/voronovsg/rocket-factory/notification/internal/converter/kafka"
	def "github.com/voronovsg/rocket-factory/notification/internal/service"
	"github.com/voronovsg/rocket-factory/platform/pkg/kafka"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
)

type service struct {
	orderAssembledConsumer kafka.Consumer
	orderAssembledDecoder  kafkaConverter.OrderAssembledDecoder
	telegramService        def.TelegramService
}

func NewService(orderAssembledConsumer kafka.Consumer, orderAssembledDecoder kafkaConverter.OrderAssembledDecoder, telegramService def.TelegramService) *service {
	return &service{
		orderAssembledConsumer: orderAssembledConsumer,
		orderAssembledDecoder:  orderAssembledDecoder,
		telegramService:        telegramService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order assembled consumer service")

	err := s.orderAssembledConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.assembled topic error", zap.Error(err))
		return err
	}

	return nil
}
