package order_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/voronovsg/rocket-factory/assembly/internal/model"
	"github.com/voronovsg/rocket-factory/platform/pkg/kafka"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
	eventsV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/events/v1"
)

type service struct {
	orderAssembledProducer kafka.Producer
}

func NewService(orderAssembledProducer kafka.Producer) *service {
	return &service{
		orderAssembledProducer: orderAssembledProducer,
	}
}

func (p *service) ProduceOrderAssembled(ctx context.Context, event model.OrderAssembledEvent) error {
	msg := &eventsV1.OrderAssembled{
		OrderUuid:    event.OrderUUID,
		UserUuid:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal OrderAssembledEvent", zap.Error(err))
		return err
	}

	err = p.orderAssembledProducer.Send(ctx, []byte(event.OrderUUID), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish OrderAssembledEvent", zap.Error(err))
		return err
	}

	return nil
}
