package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/voronovsg/rocket-factory/assembly/internal/converter/kafka"
	def "github.com/voronovsg/rocket-factory/assembly/internal/service"
	"github.com/voronovsg/rocket-factory/platform/pkg/kafka"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	orderPaidConsumer    kafka.Consumer
	orderPaidDecoder     kafkaConverter.OrderPaidDecoder
	orderProducerService def.OrderProducerService
}

func NewService(orderPaidConsumer kafka.Consumer, orderPaidDecoder kafkaConverter.OrderPaidDecoder, orderProducerService def.OrderProducerService) *service {
	return &service{
		orderPaidConsumer:    orderPaidConsumer,
		orderPaidDecoder:     orderPaidDecoder,
		orderProducerService: orderProducerService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order orderPaidConsumer service")

	err := s.orderPaidConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.paid topic error", zap.Error(err))
	}

	return err
}
