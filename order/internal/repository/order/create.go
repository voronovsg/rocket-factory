package part

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/voronovsg/rocket-factory/order/internal/model"
	repoConv "github.com/voronovsg/rocket-factory/order/internal/repository/converter"
	repoModel "github.com/voronovsg/rocket-factory/order/internal/repository/model"
)

func (r *repository) Create(_ context.Context, order model.CreateOrder) (model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	newUUID := uuid.NewString()
	newOrder := repoModel.Order{
		UUID:       newUUID,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUuids,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		CreatedAt:  time.Now(),
	}
	r.orders[newUUID] = newOrder

	return repoConv.OrderToModel(r.orders[newUUID]), nil
}
