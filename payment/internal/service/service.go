package service

import (
	"context"
)

type PaymentService interface {
	PayOrder(ctx context.Context, orderUUID, userUUID, payMethod string) (string, error)
}
