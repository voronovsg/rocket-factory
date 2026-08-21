package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/voronovsg/rocket-factory/iam/internal/model"
	repoConverter "github.com/voronovsg/rocket-factory/iam/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, newUser model.UserRegistrationInfo) (string, error) {
	newUserRepo := repoConverter.UserRegistrationToRepoModel(ctx, newUser)

	sql, args, err := sq.Insert(usersTable).
		Columns(
			userFieldLogin,
			userFieldEmail,
			userFieldPassword,
			userFieldNotificationMethods,
		).
		Values(
			newUserRepo.Login,
			newUserRepo.Email,
			newUserRepo.Password,
			newUserRepo.NotificationMethods,
		).Suffix(fmt.Sprintf("RETURNING %s", userFieldUUID)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", err
	}

	var userUUID string
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&userUUID,
	)
	if err != nil {
		return "", err
	}

	return userUUID, nil
}
