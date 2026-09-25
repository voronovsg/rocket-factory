package order

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/voronovsg/rocket-factory/order/internal/metrics"
	"github.com/voronovsg/rocket-factory/order/internal/model"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
	"github.com/voronovsg/rocket-factory/platform/pkg/ptr"
	"github.com/voronovsg/rocket-factory/platform/pkg/tracing"
)

func (s *service) PayOrder(ctx context.Context, orderUUID, paymentMethod string) (string, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		logger.Error(ctx, "Failed to get order", zap.String("orderUUID", orderUUID), zap.Error(err))
		return "", err
	}

	if order.Status != model.OrderStatusPendingPayment {
		logger.Error(ctx, "Order status is not pending payment",
			zap.String("status", order.Status),
			zap.String("orderUUID", orderUUID))
		return "", model.ErrOrderStatusInvalid
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, orderUUID, order.UserUUID, paymentMethod)
	if err != nil {
		logger.Error(ctx, "Failed to pay order", zap.String("orderUUID", orderUUID), zap.Error(err))
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return "", errors.New("payment service timeout")
		}
		return "", err
	}

	err = s.orderRepository.Update(ctx, orderUUID, model.UpdateOrder{
		TransactionUUID: &transactionUUID,
		PaymentMethod:   &paymentMethod,
		Status:          ptr.Of(model.OrderStatusPaid),
	})
	if err != nil {
		logger.Error(ctx, "Failed to update order", zap.String("orderUUID", orderUUID), zap.Error(err))
		return "", err
	}

	ctx, span := tracing.StartSpan(ctx, "produce_order_paid",
		trace.WithAttributes(
			attribute.String("order_uuid", orderUUID),
			attribute.String("user_uuid", order.UserUUID),
			attribute.String("transaction_uuid", transactionUUID),
			attribute.String("payment_method", paymentMethod),
			attribute.String("status", model.OrderStatusPaid),
			attribute.String("event_type", "order_paid"),
			attribute.String("service", "order"),
		),
		trace.WithSpanKind(trace.SpanKindProducer),
	)
	defer span.End()

	metrics.CountOrderCreated(ctx, model.OrderStatusPaid)
	metrics.SumOrderRevenue(ctx, order.TotalPrice)

	err = s.orderProducerService.ProduceOrderPaid(ctx, model.OrderPaidEvent{
		OrderUUID:       order.UUID,
		UserUUID:        order.UserUUID,
		PaymentMethod:   paymentMethod,
		TransactionUUID: transactionUUID,
		PaidAt:          time.Now(),
	})
	if err != nil {
		logger.Error(ctx, "Failed to send OrderPaidEvent", zap.String("orderUUID", orderUUID), zap.Error(err))
		return "", err
	}

	logger.Debug(ctx, "Order paid successfully",
		zap.String("orderUUID", order.UUID),
		zap.String("userUUID", order.UserUUID),
		zap.String("transactionUUID", transactionUUID),
		zap.String("paymentMethod", paymentMethod),
		zap.String("status", model.OrderStatusPaid))

	return transactionUUID, nil
}
