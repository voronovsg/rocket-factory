package part

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"

	"github.com/voronovsg/rocket-factory/inventory/internal/model"
)

func (s *ServiceSuite) TestGetPartSuccess() {
	var (
		partUuid     = gofakeit.UUID()
		expectedPart = model.Part{
			Uuid:          partUuid,
			Name:          gofakeit.Name(),
			Description:   gofakeit.LoremIpsumParagraph(1, 2, 5, ". "),
			Price:         gofakeit.Price(0, 100),
			StockQuantity: int64(gofakeit.Number(1, 100)),
			Category:      int32(gofakeit.Number(1, 2)),
		}
	)
	ctx := context.Background()

	s.mockPartRepository.On("Get", ctx, partUuid).Return(expectedPart, nil).Once()
	part, err := s.service.GetPart(ctx, partUuid)
	s.Require().NoError(err)
	s.Require().Equal(expectedPart, part)
}

func (s *ServiceSuite) TestGetPartValidationErrors() {
	testCases := []struct {
		name        string
		partUuid    string
		expectedErr error
	}{
		{
			name:        "empty part UUID",
			partUuid:    "",
			expectedErr: model.ErrPartUUIDInvalid,
		},
		{
			name:        "invalid part UUID",
			partUuid:    "invalid-part-uuid",
			expectedErr: model.ErrPartUUIDInvalid,
		},
	}
	ctx := context.Background()

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			part, err := s.service.GetPart(ctx, tc.partUuid)

			s.Require().Error(err)
			s.Require().ErrorIs(err, tc.expectedErr)
			s.Require().Empty(part)
			// валидация происходит до вызова репозитория
			s.mockPartRepository.AssertNotCalled(s.T(), "Get", mock.Anything, mock.Anything)
		})
	}
}

func (s *ServiceSuite) TestGetPartRepoErrors() {
	testCases := []struct {
		name        string
		expectedErr error
	}{
		{
			name:        "repository failure",
			expectedErr: gofakeit.Error(),
		},
		{
			name:        "part not found",
			expectedErr: model.ErrPartNotFound,
		},
	}
	ctx := context.Background()
	partUuid := gofakeit.UUID()

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.mockPartRepository.On("Get", ctx, partUuid).Return(model.Part{}, tc.expectedErr).Once()

			part, err := s.service.GetPart(ctx, partUuid)
			s.Require().Error(err)
			s.Require().ErrorIs(err, tc.expectedErr)
			s.Require().Empty(part)
		})
	}
}
