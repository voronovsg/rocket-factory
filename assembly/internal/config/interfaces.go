package config

import (
	"time"

	"github.com/IBM/sarama"
)

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

type OrderAssembledProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type OrderPaidConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

type MetricServerConfig interface {
	CollectorEndpoint() string
	CollectorInterval() time.Duration
}
