package order

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/voronovsg/rocket-factory/order/internal/model"
	"github.com/voronovsg/rocket-factory/platform/pkg/ptr"
)

func (s *ServiceSuite) TestPaySuccess() {
	var (
		orderUUID       = gofakeit.UUID()
		userUUID        = gofakeit.UUID()
		paymentMethod   = model.PaymentMethodCard
		transactionUUID = gofakeit.UUID()

		order = model.Order{
			UUID:     orderUUID,
			UserUUID: userUUID,
			Status:   model.OrderStatusPendingPayment,
		}
	)

	s.mockOrderRepository.On("Get", s.ctx, orderUUID).Return(order, nil).Once()
	s.mockPaymentClient.On("PayOrder", s.ctx, orderUUID, userUUID, paymentMethod).
		Return(transactionUUID, nil).Once()
	s.mockOrderRepository.On("Update", s.ctx, orderUUID, model.UpdateOrder{
		TransactionUUID: &transactionUUID,
		PaymentMethod:   &paymentMethod,
		Status:          ptr.Of(model.OrderStatusPaid),
	}).Return(nil).Once()

	res, err := s.service.PayOrder(s.ctx, orderUUID, paymentMethod)
	s.Require().NoError(err)
	s.Require().Equal(transactionUUID, res)
}

func (s *ServiceSuite) TestPayRepoGetError() {
	var (
		orderUUID     = gofakeit.UUID()
		paymentMethod = model.PaymentMethodCard
		repoErr       = gofakeit.Error()
	)

	s.mockOrderRepository.On("Get", s.ctx, orderUUID).Return(model.Order{}, repoErr).Once()

	_, err := s.service.PayOrder(s.ctx, orderUUID, paymentMethod)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestPayInvalidStatus() {
	var (
		orderUUID     = gofakeit.UUID()
		userUUID      = gofakeit.UUID()
		paymentMethod = model.PaymentMethodCard

		order = model.Order{
			UUID:     orderUUID,
			UserUUID: userUUID,
			Status:   model.OrderStatusPaid,
		}
	)

	s.mockOrderRepository.On("Get", s.ctx, orderUUID).Return(order, nil).Once()

	_, err := s.service.PayOrder(s.ctx, orderUUID, paymentMethod)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderStatusInvalid)
}

func (s *ServiceSuite) TestPayPaymentClientError() {
	var (
		clientErr     = gofakeit.Error()
		orderUUID     = gofakeit.UUID()
		userUUID      = gofakeit.UUID()
		paymentMethod = model.PaymentMethodCard

		order = model.Order{
			UUID:     orderUUID,
			UserUUID: userUUID,
			Status:   model.OrderStatusPendingPayment,
		}
	)

	s.mockOrderRepository.On("Get", s.ctx, orderUUID).Return(order, nil).Once()
	s.mockPaymentClient.On("PayOrder", s.ctx, orderUUID, userUUID, paymentMethod).
		Return("", clientErr).Once()

	_, err := s.service.PayOrder(s.ctx, orderUUID, paymentMethod)
	s.Require().Error(err)
	s.Require().ErrorIs(err, clientErr)
}

func (s *ServiceSuite) TestPayRepoUpdateError() {
	var (
		orderUUID       = gofakeit.UUID()
		userUUID        = gofakeit.UUID()
		paymentMethod   = model.PaymentMethodCard
		transactionUUID = gofakeit.UUID()
		repoErr         = gofakeit.Error()

		order = model.Order{
			UUID:     orderUUID,
			UserUUID: userUUID,
			Status:   model.OrderStatusPendingPayment,
		}
	)

	s.mockOrderRepository.On("Get", s.ctx, orderUUID).Return(order, nil).Once()
	s.mockPaymentClient.On("PayOrder", s.ctx, orderUUID, userUUID, paymentMethod).
		Return(transactionUUID, nil).Once()
	s.mockOrderRepository.On("Update", s.ctx, orderUUID, model.UpdateOrder{
		TransactionUUID: &transactionUUID,
		PaymentMethod:   &paymentMethod,
		Status:          ptr.Of(model.OrderStatusPaid),
	}).Return(repoErr).Once()

	_, err := s.service.PayOrder(s.ctx, orderUUID, paymentMethod)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}
