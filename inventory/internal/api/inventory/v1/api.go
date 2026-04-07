package v1

import (
	"github.com/voronovsg/rocket-factory/inventory/internal/service"
	inventoryV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer

	partService service.PartService
}

func NewAPI(partService service.PartService) *api {
	return &api{
		partService: partService,
	}
}
