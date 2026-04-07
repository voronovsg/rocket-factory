package repository

import (
	"context"

	"github.com/voronovsg/rocket-factory/inventory/internal/model"
)

// PartRepository определяет контракт для работы с деталями
type PartRepository interface {
	Get(ctx context.Context, uuid string) (model.Part, error)
	List(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
	InitGenParts(ctx context.Context) error
}
