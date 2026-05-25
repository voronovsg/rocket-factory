package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/voronovsg/rocket-factory/assembly/internal/model"
	eventsV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/events/v1"
)

type decoder struct{}

func NewOrderPaidDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.OrderPaidEvent, error) {
	var pb eventsV1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.OrderPaidEvent{
		OrderUUID:       pb.GetOrderUuid(),
		UserUUID:        pb.GetUserUuid(),
		TransactionUUID: pb.GetTransactionUuid(),
		PaidAt:          pb.GetPaidAt().AsTime(),
	}, nil
}
