package payment

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/voronovsg/rocket-factory/payment/internal/model"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
)

func (s *service) PayOrder(ctx context.Context, orderUUID, userUUID, payMethod string) (string, error) {
	_, err := uuid.Parse(orderUUID)
	if err != nil {
		return "", model.ErrOrderUUIDInvalid
	}
	_, err = uuid.Parse(userUUID)
	if err != nil {
		return "", model.ErrUserUUIDInvalid
	}

	transactionUUID := uuid.NewString()
	logger.Info(ctx, "Payment was successful",
		zap.String("orderUUID", orderUUID),
		zap.String("userUUID", userUUID),
		zap.String("transactionUUID", transactionUUID),
		zap.String("payMethod", payMethod))

	return transactionUUID, nil
}
