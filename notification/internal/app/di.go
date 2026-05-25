package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/go-telegram/bot"

	httpClient "github.com/voronovsg/rocket-factory/notification/internal/client"
	telegramClient "github.com/voronovsg/rocket-factory/notification/internal/client/http/telegram"
	"github.com/voronovsg/rocket-factory/notification/internal/config"
	kafkaConverter "github.com/voronovsg/rocket-factory/notification/internal/converter/kafka"
	"github.com/voronovsg/rocket-factory/notification/internal/converter/kafka/decoder"
	"github.com/voronovsg/rocket-factory/notification/internal/service"
	orderAssembledConsumerService "github.com/voronovsg/rocket-factory/notification/internal/service/consumer/order_assemled_consumer"
	orderPaidConsumerService "github.com/voronovsg/rocket-factory/notification/internal/service/consumer/order_paid_consumer"
	telegramService "github.com/voronovsg/rocket-factory/notification/internal/service/telegram"
	"github.com/voronovsg/rocket-factory/platform/pkg/closer"
	wrappedKafka "github.com/voronovsg/rocket-factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/voronovsg/rocket-factory/platform/pkg/kafka/consumer"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
	kafkaMiddleware "github.com/voronovsg/rocket-factory/platform/pkg/middleware/kafka"
)

type diContainer struct {
	telegramService service.TelegramService
	telegramClient  httpClient.TelegramClient
	telegramBot     *bot.Bot

	orderPaidConsumerService      service.OrderPaidConsumerService
	orderAssembledConsumerService service.OrderAssembledConsumerService

	orderPaidConsumerGroup      sarama.ConsumerGroup
	orderAssembledConsumerGroup sarama.ConsumerGroup
	orderPaidConsumer           wrappedKafka.Consumer
	orderAssembledConsumer      wrappedKafka.Consumer

	orderPaidDecoder      kafkaConverter.OrderPaidDecoder
	orderAssembledDecoder kafkaConverter.OrderAssembledDecoder
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) TelegramService() service.TelegramService {
	if d.telegramService == nil {
		d.telegramService = telegramService.NewService(d.TelegramClient())
	}

	return d.telegramService
}

func (d *diContainer) TelegramClient() httpClient.TelegramClient {
	if d.telegramClient == nil {
		d.telegramClient = telegramClient.NewClient(d.TelegramBot())
	}

	return d.telegramClient
}

func (d *diContainer) TelegramBot() *bot.Bot {
	if d.telegramBot == nil {
		b, err := bot.New(config.AppConfig().TelegramBot.Token())
		if err != nil {
			panic(fmt.Sprintf("failed to create telegram bot: %s\n", err.Error()))
		}

		d.telegramBot = b
	}

	return d.telegramBot
}

func (d *diContainer) OrderPaidConsumerService() service.OrderPaidConsumerService {
	if d.orderPaidConsumerService == nil {
		d.orderPaidConsumerService = orderPaidConsumerService.NewService(
			d.OrderPaidConsumer(),
			d.OrderPaidDecoder(),
			d.TelegramService(),
		)
	}

	return d.orderPaidConsumerService
}

func (d *diContainer) OrderAssembledConsumerService() service.OrderAssembledConsumerService {
	if d.orderAssembledConsumerService == nil {
		d.orderAssembledConsumerService = orderAssembledConsumerService.NewService(
			d.OrderAssembledConsumer(),
			d.OrderAssembledDecoder(),
			d.TelegramService(),
		)
	}

	return d.orderAssembledConsumerService
}

func (d *diContainer) OrderPaidConsumerGroup() sarama.ConsumerGroup {
	if d.orderPaidConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create order paid consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Order paid Kafka consumer group", func(ctx context.Context) error {
			return d.orderPaidConsumerGroup.Close()
		})

		d.orderPaidConsumerGroup = consumerGroup
	}

	return d.orderPaidConsumerGroup
}

func (d *diContainer) OrderAssembledConsumerGroup() sarama.ConsumerGroup {
	if d.orderAssembledConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create order assembled consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Order assembled Kafka consumer group", func(ctx context.Context) error {
			return d.orderAssembledConsumerGroup.Close()
		})

		d.orderAssembledConsumerGroup = consumerGroup
	}

	return d.orderAssembledConsumerGroup
}

func (d *diContainer) OrderPaidConsumer() wrappedKafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderPaidConsumerGroup(),
			[]string{
				config.AppConfig().OrderPaidConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderPaidConsumer
}

func (d *diContainer) OrderAssembledConsumer() wrappedKafka.Consumer {
	if d.orderAssembledConsumer == nil {
		d.orderAssembledConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderAssembledConsumerGroup(),
			[]string{
				config.AppConfig().OrderAssembledConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderAssembledConsumer
}

func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder
}

func (d *diContainer) OrderAssembledDecoder() kafkaConverter.OrderAssembledDecoder {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = decoder.NewOrderAssembledDecoder()
	}

	return d.orderAssembledDecoder
}
