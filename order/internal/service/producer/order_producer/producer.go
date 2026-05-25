package order_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/voronovsg/rocket-factory/order/internal/model"
	"github.com/voronovsg/rocket-factory/platform/pkg/kafka"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
	eventsV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/events/v1"
)

type service struct {
	orderPaidProducer kafka.Producer
}

func NewService(orderPaidProducer kafka.Producer) *service {
	return &service{
		orderPaidProducer: orderPaidProducer,
	}
}

func (p *service) ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	msg := &eventsV1.OrderPaid{
		OrderUuid:       event.OrderUUID,
		UserUuid:        event.UserUUID,
		PaymentMethod:   event.PaymentMethod,
		TransactionUuid: event.TransactionUUID,
		PaidAt:          timestamppb.New(event.PaidAt),
	}
	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "Failed to marshal order paid payload", zap.Error(err))
		return err
	}
	err = p.orderPaidProducer.Send(ctx, []byte(event.OrderUUID), payload)
	if err != nil {
		logger.Error(ctx, "Failed to send order paid payload", zap.Error(err))
	}

	return nil
}
