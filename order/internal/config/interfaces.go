package config

import "time"

type OrderHTTPConfig interface {
	Address() string
	ReadTimeout() time.Duration
}

type InventoryGRPCConfig interface {
	Address() string
}

type PaymentGRPCConfig interface {
	Address() string
}

type PostgresConfig interface {
	URI() string
	MigrationDir() string
}
