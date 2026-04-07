package order

import (
	"context"

	"github.com/voronovsg/rocket-factory/order/internal/model"
)

func (s *service) GetOrderByUUID(ctx context.Context, orderUUID string) (model.Order, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}
