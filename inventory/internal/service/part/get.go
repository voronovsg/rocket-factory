package part

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/voronovsg/rocket-factory/inventory/internal/model"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
)

func (s *service) GetPart(ctx context.Context, partUuid string) (model.Part, error) {
	_, err := uuid.Parse(partUuid)
	if err != nil {
		logger.Error(ctx, "Part UUID parse failed", zap.String("partUuid", partUuid), zap.Error(err))
		return model.Part{}, model.ErrPartUUIDInvalid
	}

	part, err := s.partRepository.Get(ctx, partUuid)
	if err != nil {
		logger.Error(ctx, "Part repository get failed", zap.String("partUuid", partUuid), zap.Error(err))
		return model.Part{}, err
	}

	logger.Debug(ctx, "Part repository get", zap.String("partUuid", partUuid))
	return part, nil
}
