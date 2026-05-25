package model

import "time"

const (
	OrderStatusAssembled      = "ASSEMBLED"
	OrderStatusPendingPayment = "PENDING_PAYMENT"
	OrderStatusPaid           = "PAID"
	OrderStatusCancelled      = "CANCELLED"

	PaymentMethodCard = "PAYMENT_METHOD_CARD"
)

type Order struct {
	UUID            string
	UserUUID        string
	PartUuids       []string
	TotalPrice      float64
	TransactionUUID *string
	PaymentMethod   *string
	Status          string
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}

type CreateOrder struct {
	UserUUID   string
	PartUuids  []string
	TotalPrice float64
	Status     string
}

type UpdateOrder struct {
	PartUuids       []string
	TotalPrice      *float64
	TransactionUUID *string
	PaymentMethod   *string
	Status          *string
}
