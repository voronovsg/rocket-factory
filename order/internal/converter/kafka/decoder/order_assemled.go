package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/voronovsg/rocket-factory/order/internal/model"
	events_v1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/events/v1"
)

type decoder struct{}

func NewOrderAssembledDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.OrderAssembledEvent, error) {
	var pb events_v1.OrderAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderAssembledEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.OrderAssembledEvent{
		OrderUUID:    pb.GetOrderUuid(),
		UserUUID:     pb.GetUserUuid(),
		BuildTimeSec: pb.GetBuildTimeSec(),
	}, nil
}
