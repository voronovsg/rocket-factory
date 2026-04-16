package payment

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/voronovsg/rocket-factory/payment/internal/model"
)

func (s *service) PayOrder(_ context.Context, orderUUID, userUUID, payMethod string) (string, error) {
	_, err := uuid.Parse(orderUUID)
	if err != nil {
		return "", model.ErrOrderUUIDInvalid
	}
	_, err = uuid.Parse(userUUID)
	if err != nil {
		return "", model.ErrUserUUIDInvalid
	}

	transactionUUID := uuid.NewString()
	log.Printf("Order UUID: %s\nUser UUID: %s\nPayment Method: %s",
		orderUUID,
		userUUID,
		payMethod)
	log.Printf("Оплата прошла успешно, Transaction UUID: %s", transactionUUID)

	return transactionUUID, nil
}
