package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/voronovsg/rocket-factory/iam/internal/model"
	commonV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/common/v1"
)

func SessionToProto(session model.Session) *commonV1.Session {
	var updatedAt *timestamppb.Timestamp
	if session.UpdatedAt != nil {
		updatedAt = timestamppb.New(*session.UpdatedAt)
	}

	return &commonV1.Session{
		Uuid:      session.UUID,
		CreatedAt: timestamppb.New(session.CreatedAt),
		UpdatedAt: updatedAt,
		ExpiresAt: timestamppb.New(session.ExpiresAt),
	}
}
