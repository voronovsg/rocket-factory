package part

import (
	"context"

	"go.uber.org/zap"

	"github.com/voronovsg/rocket-factory/inventory/internal/model"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
)

func (s *service) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	partList, err := s.partRepository.List(ctx, filter)
	if err != nil {
		logger.Error(ctx, "Part repository list", zap.Any("filter", filter), zap.Error(err))
		return nil, err
	}

	logger.Debug(ctx, "Part repository list", zap.Any("partList", partList))
	return partList, err
}
