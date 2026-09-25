package config

type PaymentGRPCConfig interface {
	Address() string
}

type PaymentHTTPConfig interface {
	Address() string
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
	EnableOTLP() bool
	CollectorEndpoint() string
	ServiceName() string
	ServiceEnvironment() string
}
