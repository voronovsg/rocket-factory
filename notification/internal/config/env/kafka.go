package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type kafkaEnvConfig struct {
	Brokers []string `env:"KAFKA_BROKERS,required"`
}

type KafkaConfig interface {
	Brokers() []string
}

type kafkaConfig struct {
	raw kafkaEnvConfig
}

func NewKafkaConfig() (KafkaConfig, error) {
	var raw kafkaEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &kafkaConfig{raw: raw}, nil
}

func (cfg *kafkaConfig) Brokers() []string {
	return cfg.raw.Brokers
}

func consumerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
