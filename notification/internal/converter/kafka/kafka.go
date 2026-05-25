package kafka

import "github.com/voronovsg/rocket-factory/notification/internal/model"

type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaidEvent, error)
}

type OrderAssembledDecoder interface {
	Decode(data []byte) (model.OrderAssembledEvent, error)
}
