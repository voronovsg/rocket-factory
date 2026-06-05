package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"github.com/voronovsg/rocket-factory/iam/internal/model"
	repoConverter "github.com/voronovsg/rocket-factory/iam/internal/repository/converter"
	repoModel "github.com/voronovsg/rocket-factory/iam/internal/repository/model"
)

func (r *repository) GetByIdentifier(ctx context.Context, identifier model.UserIdentifier) (model.User, error) {
	q := sq.Select(
		userFieldUUID,
		userFieldLogin,
		userFieldEmail,
		userFieldNotificationMethods,
		userFieldPassword,
		userFieldCreatedAt,
		userFieldUpdatedAt,
	).
		From(usersTable).
		PlaceholderFormat(sq.Dollar).
		Limit(1)

	// Только одно из полей должно быть задано
	switch {
	case identifier.UUID != nil:
		q = q.Where(sq.Eq{userFieldUUID: *identifier.UUID})
	case identifier.Login != nil:
		q = q.Where(sq.Eq{userFieldLogin: *identifier.Login})
	case identifier.Email != nil:
		q = q.Where(sq.Eq{userFieldEmail: *identifier.Email})
	default:
		return model.User{}, model.ErrUserIdentifierInvalid
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return model.User{}, err
	}

	var user repoModel.User
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&user.UUID,
		&user.Login,
		&user.Email,
		&user.NotificationMethods,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, model.ErrUserNotFound
		}

		return model.User{}, err
	}

	return repoConverter.UserToModel(ctx, user), nil
}
