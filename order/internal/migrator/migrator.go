package migrator

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

type migrator struct {
	db             *sql.DB
	migrationsPath string
}

func New(db *sql.DB, migrationsPath string) *migrator {
	return &migrator{
		db:             db,
		migrationsPath: migrationsPath,
	}
}

func (m *migrator) Up() error {
	err := goose.Up(m.db, m.migrationsPath)
	if err != nil {
		return err
	}
	return nil
}
