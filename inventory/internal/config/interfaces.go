package config

type InventoryGRPCConfig interface {
	Address() string
}

type InventoryHTTPConfig interface {
	Address() string
}

type MongoConfig interface {
	URI() string
	DBName() string
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type IAMGRPCConfig interface {
	Address() string
}
