package kafka

import "github.com/voronovsg/rocket-factory/assembly/internal/model"

type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaidEvent, error)
}
