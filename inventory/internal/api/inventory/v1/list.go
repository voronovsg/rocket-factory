package v1

import (
	"context"

	"github.com/voronovsg/rocket-factory/inventory/internal/converter"
	inventoryV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts, err := a.partService.ListParts(ctx, converter.PartsFilterToModel(req.GetFilter()))
	if err != nil {
		return nil, err
	}

	return &inventoryV1.ListPartsResponse{
		Parts: converter.PartListToProto(parts),
	}, nil
}
