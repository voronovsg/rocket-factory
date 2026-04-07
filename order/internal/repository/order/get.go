package part

import (
	"context"

	"github.com/voronovsg/rocket-factory/order/internal/model"
	repoConv "github.com/voronovsg/rocket-factory/order/internal/repository/converter"
)

func (r *repository) Get(_ context.Context, orderUUID string) (model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	order, ok := r.orders[orderUUID]
	if !ok {
		return model.Order{}, model.ErrOrderNotFound
	}

	return repoConv.OrderToModel(order), nil
}
