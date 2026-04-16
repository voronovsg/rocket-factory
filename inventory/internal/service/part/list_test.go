package part

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/voronovsg/rocket-factory/inventory/internal/model"
)

func (s *ServiceSuite) TestListPartsSuccess() {
	var (
		part1    = gofakeit.UUID()
		part2    = gofakeit.UUID()
		part3    = gofakeit.UUID()
		name     = gofakeit.Name()
		category = int32(gofakeit.Number(0, 3))

		filter = model.PartsFilter{
			Uuids:      []string{part1, part2, part3},
			Names:      []string{name},
			Categories: []int32{category},
		}

		expectedPartList = []model.Part{
			{
				Uuid:          part1,
				Name:          name,
				Description:   gofakeit.LoremIpsumParagraph(1, 2, 5, ". "),
				Price:         gofakeit.Price(0, 100),
				StockQuantity: int64(gofakeit.Number(0, 100)),
				Category:      category,
			},
			{
				Uuid:          part2,
				Name:          name,
				Description:   gofakeit.LoremIpsumParagraph(1, 2, 5, ". "),
				Price:         gofakeit.Price(0, 100),
				StockQuantity: int64(gofakeit.Number(0, 100)),
				Category:      category,
			},
		}
	)
	ctx := context.Background()

	s.mockPartRepository.On("List", ctx, filter).Return(expectedPartList, nil).Once()
	partList, err := s.service.ListParts(ctx, filter)
	s.Require().NoError(err)
	s.Require().Equal(expectedPartList, partList)
}

func (s *ServiceSuite) TestListPartsRepoErrors() {
	var (
		part1       = gofakeit.UUID()
		part2       = gofakeit.UUID()
		part3       = gofakeit.UUID()
		name        = gofakeit.Name()
		category    = int32(gofakeit.Number(0, 3))
		expectedErr = gofakeit.Error()

		filter = model.PartsFilter{
			Uuids:      []string{part1, part2, part3},
			Names:      []string{name},
			Categories: []int32{category},
		}
	)
	ctx := context.Background()

	s.mockPartRepository.On("List", ctx, filter).Return(nil, expectedErr).Once()
	_, err := s.service.ListParts(ctx, filter)
	s.Require().Error(err)
	s.Require().ErrorIs(err, expectedErr)
}
