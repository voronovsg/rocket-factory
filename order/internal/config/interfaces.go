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
}

type PaymentGRPCConfig interface {
	Address() string
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
