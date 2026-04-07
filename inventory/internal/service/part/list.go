package part

import (
	"context"

	"github.com/voronovsg/rocket-factory/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	partList, err := s.partRepository.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return partList, err
}
