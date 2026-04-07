package v1

import (
	"context"

	clientConv "github.com/voronovsg/rocket-factory/order/internal/client/converter"
	"github.com/voronovsg/rocket-factory/order/internal/model"
	generatedInventoryV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	parts, err := c.generatedClient.ListParts(ctx, &generatedInventoryV1.ListPartsRequest{
		Filter: clientConv.PartsFilterToProto(filter),
	})
	if err != nil {
		return nil, err
	}

	return clientConv.PartListToModel(parts.Parts), nil
}
