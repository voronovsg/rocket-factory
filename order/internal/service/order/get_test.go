package order

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/voronovsg/rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestGetSuccess() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()

		expectedOrder = model.Order{
			UUID:       orderUUID,
			UserUUID:   userUUID,
			PartUuids:  []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice: gofakeit.Price(10, 100),
			Status:     model.OrderStatusPendingPayment,
			CreatedAt:  gofakeit.Date(),
		}
	)

	s.mockOrderRepository.On("Get", s.ctx, orderUUID).Return(expectedOrder, nil).Once()

	order, err := s.service.GetOrderByUUID(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Require().Equal(expectedOrder, order)
}

func (s *ServiceSuite) TestGetRepoError() {
	var (
		repoErr   = gofakeit.Error()
		orderUUID = gofakeit.UUID()
	)

	s.mockOrderRepository.On("Get", s.ctx, orderUUID).Return(model.Order{}, repoErr).Once()

	_, err := s.service.GetOrderByUUID(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestGetNotFound() {
	orderUUID := gofakeit.UUID()

	s.mockOrderRepository.On("Get", s.ctx, orderUUID).Return(model.Order{}, model.ErrOrderNotFound).Once()

	_, err := s.service.GetOrderByUUID(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderNotFound)
}
