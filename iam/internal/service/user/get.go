package user

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/voronovsg/rocket-factory/iam/internal/model"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
)

func (s *service) Get(ctx context.Context, userUUID string) (model.User, error) {
	_, err := uuid.Parse(userUUID)
	if err != nil {
		return model.User{}, model.ErrUserUUIDInvalid
	}

	user, err := s.userRepository.GetByIdentifier(ctx, model.UserIdentifier{UUID: &userUUID})
	if err != nil {
		logger.Error(ctx, "failed to get user", zap.Error(err))
		return model.User{}, err
	}

	return user, nil
}
