package user

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/voronovsg/rocket-factory/iam/internal/repository"
)

var _ def.UserRepository = (*repository)(nil)

const (
	usersTable = "users"

	userFieldUUID                = "uuid"
	userFieldLogin               = "login"
	userFieldEmail               = "email"
	userFieldPassword            = "password"
	userFieldNotificationMethods = "notification_methods"
	userFieldCreatedAt           = "created_at"
	userFieldUpdatedAt           = "updated_at"
)

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *repository {
	return &repository{db: db}
}
