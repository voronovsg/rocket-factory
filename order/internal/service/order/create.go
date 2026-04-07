package order

import (
	"context"
	"log"

	"github.com/go-faster/errors"

	"github.com/voronovsg/rocket-factory/order/internal/model"
)

func (s *service) CreateOrder(ctx context.Context, createOrder model.CreateOrder) (model.Order, error) {
	parts, err := s.inventoryClient.ListParts(ctx, model.PartsFilter{
		Uuids: createOrder.PartUuids,
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return model.Order{}, errors.New("inventory service timeout")
		}

		return model.Order{}, errors.Errorf("inventory service error: %s", err.Error())
	}

	if len(parts) != len(createOrder.PartUuids) {
		return model.Order{}, model.ErrPartsNotFound
	}

	var totalPrice float64
	for _, part := range parts {
		totalPrice += part.Price
	}

	createOrder.TotalPrice = totalPrice
	createOrder.Status = model.OrderStatusPendingPayment

	order, err := s.orderRepository.Create(ctx, createOrder)
	if err != nil {
		return model.Order{}, err
	}

	log.Printf(`[Order Created]
	Order UUID: %s
	User UUID: %s
	Part UUIDs: %v
	Total Price: %f
	Status: %s
	CreatedAt: %v`, order.UUID, order.UserUUID, order.PartUuids, order.TotalPrice, order.Status, order.CreatedAt)

	return order, nil
}
