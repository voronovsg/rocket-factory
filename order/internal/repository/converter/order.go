package converter

import (
	"github.com/voronovsg/rocket-factory/order/internal/model"
	repoModel "github.com/voronovsg/rocket-factory/order/internal/repository/model"
)

func OrderToModel(order repoModel.Order) model.Order {
	return model.Order{
		UUID:            order.UUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   order.PaymentMethod,
		Status:          order.Status,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
}
