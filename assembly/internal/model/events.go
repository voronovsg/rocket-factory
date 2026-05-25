package model

import (
	"time"
)

type OrderPaidEvent struct {
	OrderUUID       string
	UserUUID        string
	TransactionUUID string
	PaymentMethod   string
	PaidAt          time.Time
}

type OrderAssembledEvent struct {
	OrderUUID    string
	UserUUID     string
	BuildTimeSec int64
}
