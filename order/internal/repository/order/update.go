package part

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/voronovsg/rocket-factory/order/internal/model"
)

func (r *repository) Update(ctx context.Context, orderUUID string, updateOrder model.UpdateOrder) error {
	query := sq.Update(ordersTable).
		Set(ordersFieldUpdatedAt, sq.Expr("now()")).
		Where(sq.Eq{ordersFieldUUID: orderUUID}).
		PlaceholderFormat(sq.Dollar)

	if updateOrder.PartUuids != nil {
		query = query.Set(ordersFieldPartUuids, updateOrder.PartUuids)
	}
	if updateOrder.TotalPrice != nil {
		query = query.Set(ordersFieldTotalPrice, *updateOrder.TotalPrice)
	}
	if updateOrder.TransactionUUID != nil {
		query = query.Set(ordersFieldTransactionUUID, *updateOrder.TransactionUUID)
	}
	if updateOrder.PaymentMethod != nil {
		query = query.Set(ordersFieldPaymentMethod, *updateOrder.PaymentMethod)
	}
	if updateOrder.Status != nil {
		query = query.Set(ordersFieldStatus, *updateOrder.Status)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return model.ErrOrderNotFound
	}

	return nil
}
