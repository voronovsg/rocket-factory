package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/voronovsg/rocket-factory/notification/internal/config/env"
)

var appConfig *config

type config struct {
	Logger                 LoggerConfig
	TelegramBot            TelegramBotConfig
	Kafka                  KafkaConfig
	OrderPaidConsumer      OrderPaidConsumerConfig
	OrderAssembledConsumer OrderAssembledConsumerConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	telegramBotCfg, err := env.NewTelegramBotConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderPaidConsumerCfg, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}

	orderAssembledConsumerCfg, err := env.NewOrderAssembledConsumerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                 loggerCfg,
		TelegramBot:            telegramBotCfg,
		Kafka:                  kafkaCfg,
		OrderPaidConsumer:      orderPaidConsumerCfg,
		OrderAssembledConsumer: orderAssembledConsumerCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
