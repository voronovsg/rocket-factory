package config

import (
	"time"

	"github.com/IBM/sarama"
)

type OrderHTTPConfig interface {
	Address() string
	ReadTimeout() time.Duration
}

type InventoryGRPCConfig interface {
	Address() string
	ServiceName() string
}

type PaymentGRPCConfig interface {
	Address() string
	ServiceName() string
}

type IAMGRPCConfig interface {
	Address() string
}

type PostgresConfig interface {
	URI() string
	MigrationDir() string
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
	EnableOTLP() bool
	CollectorEndpoint() string
	ServiceName() string
	ServiceEnvironment() string
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type OrderAssembledConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

type TraceConfig interface {
	CollectorEndpoint() string
	ServiceName() string
	Environment() string
	ServiceVersion() string
}

type MetricServerConfig interface {
	CollectorEndpoint() string
	CollectorInterval() time.Duration
}
