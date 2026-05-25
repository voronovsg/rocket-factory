package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/voronovsg/rocket-factory/order/internal/converter/kafka"
	"github.com/voronovsg/rocket-factory/order/internal/repository"
	def "github.com/voronovsg/rocket-factory/order/internal/service"
	"github.com/voronovsg/rocket-factory/platform/pkg/kafka"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	orderAssembledConsumer kafka.Consumer
	orderAssembledDecoder  kafkaConverter.OrderAssembledDecoder
	orderRepository        repository.OrderRepository
}

func NewService(
	orderAssembledConsumer kafka.Consumer,
	orderAssembledDecoder kafkaConverter.OrderAssembledDecoder,
	orderRepository repository.OrderRepository,
) *service {
	return &service{
		orderAssembledConsumer: orderAssembledConsumer,
		orderAssembledDecoder:  orderAssembledDecoder,
		orderRepository:        orderRepository,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting orderAssembledConsumer")
	err := s.orderAssembledConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Consumption from topic order.assembled error", zap.Error(err))
		return err
	}

	return nil
}
