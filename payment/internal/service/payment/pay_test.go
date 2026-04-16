package payment

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/voronovsg/rocket-factory/payment/internal/model"
)

func (s *ServiceSuite) TestPayOrderSuccess() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()
		payMethod = "PAYMENT_METHOD_CARD"
	)
	ctx := context.Background()

	res, err := s.service.PayOrder(ctx, orderUUID, userUUID, payMethod)
	s.Require().NoError(err)
	s.Require().NotEmpty(res)
	_, parseErr := uuid.Parse(res)
	s.Require().NoError(parseErr)
}

func (s *ServiceSuite) TestPayOrderValidationErrors() {
	testCases := []struct {
		name        string
		orderUUID   string
		userUUID    string
		payMethod   string
		expectedErr error
	}{
		{
			name:        "empty order UUID",
			orderUUID:   "",
			userUUID:    gofakeit.UUID(),
			payMethod:   "PAYMENT_METHOD_SBP",
			expectedErr: model.ErrOrderUUIDInvalid,
		},
		{
			name:        "invalid order UUID",
			orderUUID:   "invalid-order-uuid",
			userUUID:    gofakeit.UUID(),
			payMethod:   "PAYMENT_METHOD_SBP",
			expectedErr: model.ErrOrderUUIDInvalid,
		},
		{
			name:        "empty user UUID",
			orderUUID:   gofakeit.UUID(),
			userUUID:    "",
			payMethod:   "PAYMENT_METHOD_SBP",
			expectedErr: model.ErrUserUUIDInvalid,
		},
		{
			name:        "invalid user UUID",
			orderUUID:   gofakeit.UUID(),
			userUUID:    "invalid-user-uuid",
			payMethod:   "PAYMENT_METHOD_SBP",
			expectedErr: model.ErrUserUUIDInvalid,
		},
	}
	ctx := context.Background()

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			res, err := s.service.PayOrder(ctx, tc.orderUUID, tc.userUUID, tc.payMethod)
			s.Require().Error(err)
			s.Require().ErrorIs(err, tc.expectedErr)
			s.Require().Empty(res)
		})
	}
}
