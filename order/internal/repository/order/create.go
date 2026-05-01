package part

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/voronovsg/rocket-factory/order/internal/model"
	repoConv "github.com/voronovsg/rocket-factory/order/internal/repository/converter"
	repoModel "github.com/voronovsg/rocket-factory/order/internal/repository/model"
)

func (r *repository) Create(ctx context.Context, order model.CreateOrder) (model.Order, error) {
	sql, args, err := sq.Insert(ordersTable).
		Columns(
			ordersFieldUserUUID,
			ordersFieldPartUuids,
			ordersFieldTotalPrice,
			ordersFieldStatus,
		).
		Values(
			order.UserUUID,
			order.PartUuids,
			order.TotalPrice,
			order.Status,
		).
		Suffix(
			fmt.Sprintf("RETURNING %s, %s, %s, %s, %s, %s, %s, %s, %s",
				ordersFieldUUID,
				ordersFieldUserUUID,
				ordersFieldPartUuids,
				ordersFieldTotalPrice,
				ordersFieldTransactionUUID,
				ordersFieldPaymentMethod,
				ordersFieldStatus,
				ordersFieldCreatedAt,
				ordersFieldUpdatedAt,
			)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return model.Order{}, err
	}

	var newOrder repoModel.Order
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&newOrder.UUID,
		&newOrder.UserUUID,
		&newOrder.PartUuids,
		&newOrder.TotalPrice,
		&newOrder.TransactionUUID,
		&newOrder.PaymentMethod,
		&newOrder.Status,
		&newOrder.CreatedAt,
		&newOrder.UpdatedAt,
	)
	if err != nil {
		return model.Order{}, err
	}

	return repoConv.OrderToModel(newOrder), nil
}
