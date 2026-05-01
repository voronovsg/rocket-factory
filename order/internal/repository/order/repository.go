package part

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/voronovsg/rocket-factory/order/internal/repository"
)

const (
	ordersTable                = "orders"
	ordersFieldUUID            = "uuid"
	ordersFieldUserUUID        = "user_uuid"
	ordersFieldPartUuids       = "part_uuids"
	ordersFieldTotalPrice      = "total_price"
	ordersFieldTransactionUUID = "transaction_uuid"
	ordersFieldPaymentMethod   = "payment_method"
	ordersFieldStatus          = "status"
	ordersFieldCreatedAt       = "created_at"
	ordersFieldUpdatedAt       = "updated_at"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		db: pool,
	}
}
