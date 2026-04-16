package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/voronovsg/rocket-factory/inventory/internal/model"
)

func (s *service) GetPart(ctx context.Context, partUuid string) (model.Part, error) {
	_, err := uuid.Parse(partUuid)
	if err != nil {
		return model.Part{}, model.ErrPartUUIDInvalid
	}

	part, err := s.partRepository.Get(ctx, partUuid)
	if err != nil {
		return model.Part{}, err
	}

	return part, nil
}
