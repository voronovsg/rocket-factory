package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/voronovsg/rocket-factory/notification/internal/model"
	eventsV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/events/v1"
)

type orderAssembledDecoder struct{}

func NewOrderAssembledDecoder() *orderAssembledDecoder {
	return &orderAssembledDecoder{}
}

func (d *orderAssembledDecoder) Decode(data []byte) (model.OrderAssembledEvent, error) {
	var pb eventsV1.OrderAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderAssembledEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.OrderAssembledEvent{
		OrderUUID:    pb.GetOrderUuid(),
		BuildTimeSec: pb.GetBuildTimeSec(),
	}, nil
}
