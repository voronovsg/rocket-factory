package order

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/voronovsg/rocket-factory/order/internal/model"
	"github.com/voronovsg/rocket-factory/shared/pkg/ptr"
)

func (s *ServiceSuite) TestCancelOrderByUUIDSuccess() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()
		order     = model.Order{
			UUID:     orderUUID,
			UserUUID: userUUID,
			Status:   model.OrderStatusPendingPayment,
		}
	)

	s.mockOrderRepository.On("Get", s.ctx, orderUUID).Return(order, nil).Once()
	s.mockOrderRepository.On("Update", s.ctx, orderUUID, model.UpdateOrder{
		Status: ptr.Of(model.OrderStatusCancelled),
	}).Return(nil).Once()

	err := s.service.CancelOrderByUUID(s.ctx, orderUUID)
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestCancelOrderByUUIDInvalidStatus() {
	var (
		orderUUID = gofakeit.UUID()
		order     = model.Order{
			UUID:     orderUUID,
			UserUUID: gofakeit.UUID(),
			Status:   model.OrderStatusPaid,
		}
	)

	s.mockOrderRepository.On("Get", s.ctx, orderUUID).Return(order, nil)
	err := s.service.CancelOrderByUUID(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderStatusInvalid)
}

func (s *ServiceSuite) TestCancelOrderByUUIDRepoErrors() {
	var (
		expectedErr = gofakeit.Error()
		orderUUID   = gofakeit.UUID()
	)

	s.mockOrderRepository.On("Get", s.ctx, orderUUID).Return(model.Order{}, expectedErr)
	err := s.service.CancelOrderByUUID(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, expectedErr)
}

func (s *ServiceSuite) TestCancelOrderByUUIDUpdateError() {
	var (
		orderUUID   = gofakeit.UUID()
		userUUID    = gofakeit.UUID()
		expectedErr = gofakeit.Error()
		order       = model.Order{
			UUID:     orderUUID,
			UserUUID: userUUID,
			Status:   model.OrderStatusPendingPayment,
		}
	)

	s.mockOrderRepository.On("Get", s.ctx, orderUUID).Return(order, nil)
	s.mockOrderRepository.On("Update", s.ctx, orderUUID, model.UpdateOrder{
		Status: ptr.Of(model.OrderStatusCancelled),
	}).Return(expectedErr)

	err := s.service.CancelOrderByUUID(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, expectedErr)
}
