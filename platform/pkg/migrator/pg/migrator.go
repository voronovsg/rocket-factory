package pg

import (
	"database/sql"

	"github.com/pressly/goose/v3"

	def "github.com/voronovsg/rocket-factory/platform/pkg/migrator"
)

var _ def.Migrator = (*migrator)(nil)

type migrator struct {
	db            *sql.DB
	migrationsDir string
}

func NewMigrator(db *sql.DB, migrationsDir string) *migrator {
	return &migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

func (m *migrator) Up() error {
	err := goose.Up(m.db, m.migrationsDir)
	if err != nil {
		return err
	}

	return nil
}
