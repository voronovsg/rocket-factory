package part

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"github.com/voronovsg/rocket-factory/order/internal/model"
	repoConv "github.com/voronovsg/rocket-factory/order/internal/repository/converter"
	repoModel "github.com/voronovsg/rocket-factory/order/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, orderUUID string) (model.Order, error) {
	sql, args, err := sq.Select(
		ordersFieldUUID,
		ordersFieldUserUUID,
		ordersFieldPartUuids,
		ordersFieldTotalPrice,
		ordersFieldTransactionUUID,
		ordersFieldPaymentMethod,
		ordersFieldStatus,
		ordersFieldCreatedAt,
		ordersFieldUpdatedAt,
	).
		From(ordersTable).
		Where(sq.Eq{ordersFieldUUID: orderUUID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return model.Order{}, err
	}

	var order repoModel.Order
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&order.UUID,
		&order.UserUUID,
		&order.PartUuids,
		&order.TotalPrice,
		&order.TransactionUUID,
		&order.PaymentMethod,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Order{}, model.ErrOrderNotFound
		}
		return model.Order{}, err
	}

	return repoConv.OrderToModel(order), nil
}
