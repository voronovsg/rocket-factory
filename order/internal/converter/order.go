package converter

import (
	"github.com/google/uuid"

	"github.com/voronovsg/rocket-factory/order/internal/model"
	orderV1 "github.com/voronovsg/rocket-factory/shared/pkg/openapi/order/v1"
)

func CreateOrderToModel(createOrder *orderV1.CreateOrderRequest) model.CreateOrder {
	partUuids := make([]string, 0, len(createOrder.PartUuids))
	for _, partUuid := range createOrder.PartUuids {
		partUuids = append(partUuids, partUuid.String())
	}

	return model.CreateOrder{
		UserUUID:  createOrder.UserUUID.String(),
		PartUuids: partUuids,
	}
}

func OrderToProto(order model.Order) orderV1.OrderDto {
	partUUIDs := make([]uuid.UUID, 0, len(order.PartUuids))
	for _, partUUID := range order.PartUuids {
		partUUIDs = append(partUUIDs, uuid.MustParse(partUUID))
	}

	var transactionUUID orderV1.OptNilUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderV1.NewOptNilUUID(uuid.MustParse(*order.TransactionUUID))
	}

	var paymentMethod orderV1.OptPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderV1.NewOptPaymentMethod(orderV1.PaymentMethod(*order.PaymentMethod))
	}

	var updatedAt orderV1.OptDateTime
	if order.UpdatedAt != nil {
		updatedAt = orderV1.NewOptDateTime(*order.UpdatedAt)
	}

	return orderV1.OrderDto{
		UUID:            uuid.MustParse(order.UUID),
		UserUUID:        uuid.MustParse(order.UserUUID),
		PartUuids:       partUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderV1.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       updatedAt,
	}
}
