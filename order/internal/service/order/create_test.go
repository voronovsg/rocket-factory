package order

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"

	"github.com/voronovsg/rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestCreateOrderSuccess() {
	var (
		userUUID  = gofakeit.UUID()
		partUUID1 = gofakeit.UUID()
		partUUID2 = gofakeit.UUID()
		partUuids = []string{partUUID1, partUUID2}
		createdAt = gofakeit.Date()

		part1 = model.Part{
			Uuid:  partUUID1,
			Price: gofakeit.Price(0, 100),
		}
		part2 = model.Part{
			Uuid:  partUUID2,
			Price: gofakeit.Price(0, 100),
		}
		parts = []model.Part{part1, part2}

		totalPrice = part1.Price + part2.Price

		createOrder = model.CreateOrder{
			UserUUID:  userUUID,
			PartUuids: partUuids,
		}

		expectedOrder = model.Order{
			UUID:       gofakeit.UUID(),
			UserUUID:   userUUID,
			PartUuids:  partUuids,
			TotalPrice: totalPrice,
			Status:     model.OrderStatusPendingPayment,
			CreatedAt:  createdAt,
		}
	)

	s.mockInventoryClient.On("ListParts", s.ctx, model.PartsFilter{
		Uuids: partUuids,
	}).Return(parts, nil).Once()

	s.mockOrderRepository.On("Create", s.ctx, mock.MatchedBy(func(order model.CreateOrder) bool {
		s.Require().Equal(expectedOrder.UserUUID, order.UserUUID)
		s.Require().Equal(expectedOrder.PartUuids, order.PartUuids)
		s.Require().Equal(expectedOrder.TotalPrice, order.TotalPrice)
		s.Require().Equal(expectedOrder.Status, order.Status)
		return true
	})).Return(expectedOrder, nil).Once()

	order, err := s.service.CreateOrder(s.ctx, createOrder)
	s.Require().NoError(err)
	s.Require().NotEmpty(order.UUID)
	s.Require().Equal(order.UserUUID, userUUID)
	s.Require().Equal(order.PartUuids, partUuids)
	s.Require().Equal(order.TotalPrice, totalPrice)
	s.Require().Equal(order.Status, model.OrderStatusPendingPayment)
	s.Require().NotEmpty(order.CreatedAt)
}

func (s *ServiceSuite) TestCreateInventoryServiceTimeout() {
	var (
		userUUID  = gofakeit.UUID()
		partUUID1 = gofakeit.UUID()
		partUUID2 = gofakeit.UUID()
		partUuids = []string{partUUID1, partUUID2}

		createOrder = model.CreateOrder{
			UserUUID:  userUUID,
			PartUuids: partUuids,
		}
	)

	s.mockInventoryClient.On("ListParts", s.ctx, model.PartsFilter{
		Uuids: partUuids,
	}).Return([]model.Part{}, context.DeadlineExceeded).Once()

	_, err := s.service.CreateOrder(s.ctx, createOrder)
	s.Require().Error(err)
	s.Require().ErrorContains(err, "inventory service timeout")
}

func (s *ServiceSuite) TestCreateInventoryServiceError() {
	var (
		userUUID  = gofakeit.UUID()
		partUUID1 = gofakeit.UUID()
		partUUID2 = gofakeit.UUID()
		partUuids = []string{partUUID1, partUUID2}

		customErr = gofakeit.Error()

		createOrder = model.CreateOrder{
			UserUUID:  userUUID,
			PartUuids: partUuids,
		}
	)

	s.mockInventoryClient.On("ListParts", s.ctx, model.PartsFilter{
		Uuids: partUuids,
	}).Return([]model.Part{}, customErr).Once()

	_, err := s.service.CreateOrder(s.ctx, createOrder)
	s.Require().Error(err)
	s.Require().ErrorContains(err, "inventory service error")
	s.Require().ErrorContains(err, customErr.Error())
}

func (s *ServiceSuite) TestCreateSomePartsNotFound() {
	var (
		userUUID  = gofakeit.UUID()
		partUUID1 = gofakeit.UUID()
		partUUID2 = gofakeit.UUID()
		partUuids = []string{partUUID1, partUUID2}

		part1 = model.Part{
			Uuid:  partUUID1,
			Price: gofakeit.Price(0, 100),
		}

		parts = []model.Part{part1}

		createOrder = model.CreateOrder{
			UserUUID:  userUUID,
			PartUuids: partUuids,
		}
	)

	s.mockInventoryClient.On("ListParts", s.ctx, model.PartsFilter{
		Uuids: partUuids,
	}).Return(parts, nil).Once()

	_, err := s.service.CreateOrder(s.ctx, createOrder)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrPartsNotFound)
}

func (s *ServiceSuite) TestCreateRepoError() {
	var (
		userUUID  = gofakeit.UUID()
		partUUID1 = gofakeit.UUID()
		partUUID2 = gofakeit.UUID()
		partUuids = []string{partUUID1, partUUID2}

		repoErr = gofakeit.Error()

		part1 = model.Part{
			Uuid:  partUUID1,
			Price: gofakeit.Price(0, 100),
		}
		part2 = model.Part{
			Uuid:  partUUID2,
			Price: gofakeit.Price(0, 100),
		}
		parts = []model.Part{part1, part2}

		totalPrice = part1.Price + part2.Price

		createOrder = model.CreateOrder{
			UserUUID:  userUUID,
			PartUuids: partUuids,
		}
	)

	s.mockInventoryClient.On("ListParts", s.ctx, model.PartsFilter{
		Uuids: partUuids,
	}).Return(parts, nil).Once()

	s.mockOrderRepository.On("Create", s.ctx, mock.MatchedBy(func(order model.CreateOrder) bool {
		s.Require().Equal(userUUID, order.UserUUID)
		s.Require().Equal(partUuids, order.PartUuids)
		s.Require().Equal(totalPrice, order.TotalPrice)
		s.Require().Equal(model.OrderStatusPendingPayment, order.Status)

		return true
	})).Return(model.Order{}, repoErr).Once()

	_, err := s.service.CreateOrder(s.ctx, createOrder)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}
