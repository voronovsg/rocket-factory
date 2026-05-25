package service

import (
	"context"

	"github.com/voronovsg/rocket-factory/notification/internal/model"
)

type OrderPaidConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderAssembledConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type TelegramService interface {
	SendOrderPaidNotification(ctx context.Context, event model.OrderPaidEvent) error
	SendOrderAssembledNotification(ctx context.Context, event model.OrderAssembledEvent) error
}
