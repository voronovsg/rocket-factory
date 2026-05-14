package migrator

type Migrator interface {
	Up() error
}
