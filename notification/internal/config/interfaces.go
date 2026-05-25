package config

import "github.com/IBM/sarama"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type TelegramBotConfig interface {
	Token() string
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

type OrderAssembledConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}
