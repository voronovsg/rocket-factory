package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1Impl "github.com/voronovsg/rocket-factory/order/internal/api/order/v1"
	grpcClient "github.com/voronovsg/rocket-factory/order/internal/client/grpc"
	inventoryV1 "github.com/voronovsg/rocket-factory/order/internal/client/grpc/inventory/v1"
	paymentV1 "github.com/voronovsg/rocket-factory/order/internal/client/grpc/payment/v1"
	"github.com/voronovsg/rocket-factory/order/internal/config"
	kafkaConverter "github.com/voronovsg/rocket-factory/order/internal/converter/kafka"
	"github.com/voronovsg/rocket-factory/order/internal/converter/kafka/decoder"
	"github.com/voronovsg/rocket-factory/order/internal/repository"
	orderRepository "github.com/voronovsg/rocket-factory/order/internal/repository/order"
	"github.com/voronovsg/rocket-factory/order/internal/service"
	orderConsumerService "github.com/voronovsg/rocket-factory/order/internal/service/consumer/order_consumer"
	orderService "github.com/voronovsg/rocket-factory/order/internal/service/order"
	orderProducerService "github.com/voronovsg/rocket-factory/order/internal/service/producer/order_producer"
	"github.com/voronovsg/rocket-factory/platform/pkg/closer"
	wrapKafka "github.com/voronovsg/rocket-factory/platform/pkg/kafka"
	wrapKafkaConsumer "github.com/voronovsg/rocket-factory/platform/pkg/kafka/consumer"
	wrapKafkaProducer "github.com/voronovsg/rocket-factory/platform/pkg/kafka/producer"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
	kafkaMiddleware "github.com/voronovsg/rocket-factory/platform/pkg/middleware/kafka"
	"github.com/voronovsg/rocket-factory/platform/pkg/migrator"
	"github.com/voronovsg/rocket-factory/platform/pkg/migrator/pg"
	orderV1 "github.com/voronovsg/rocket-factory/shared/pkg/openapi/order/v1"
	generatedInventoryV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/inventory/v1"
	generatedPaymentV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderV1API           orderV1.Handler
	orderService         service.OrderService
	orderRepository      repository.OrderRepository
	orderProducerService service.OrderProducerService
	orderConsumerService service.ConsumerService

	inventoryClient grpcClient.InventoryClient
	paymentClient   grpcClient.PaymentClient
	migrationRunner migrator.Migrator

	pgPool      *pgxpool.Pool
	pgPoolCfg   *pgxpool.Config
	stdLibPgCon *sql.DB

	consumerGroup          sarama.ConsumerGroup
	orderAssembledConsumer wrapKafka.Consumer

	orderAssembledDecoder kafkaConverter.OrderAssembledDecoder
	syncProducer          sarama.SyncProducer
	orderPaidProducer     wrapKafka.Producer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderV1API(ctx context.Context) orderV1.Handler {
	if d.orderV1API == nil {
		d.orderV1API = orderV1Impl.NewAPI(d.OrderService(ctx))
	}

	return d.orderV1API
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewService(
			d.OrderRepository(ctx),
			d.InventoryClient(),
			d.PaymentClient(),
			d.OrderProducerService(),
		)
	}

	return d.orderService
}

func (d *diContainer) OrderProducerService() service.OrderProducerService {
	if d.orderProducerService == nil {
		d.orderProducerService = orderProducerService.NewService(
			d.OrderPaidProducer(),
		)
	}

	return d.orderProducerService
}

func (d *diContainer) OrderConsumerService(ctx context.Context) service.ConsumerService {
	if d.orderConsumerService == nil {
		d.orderConsumerService = orderConsumerService.NewService(
			d.OrderAssembledConsumer(),
			d.OrderAssembledDecoder(),
			d.OrderRepository(ctx),
		)
	}

	return d.orderConsumerService
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewRepository(d.PgPool(ctx))
	}

	return d.orderRepository
}

func (d *diContainer) InventoryClient() grpcClient.InventoryClient {
	if d.inventoryClient == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().InventoryGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("❌ не удалось подключиться к grpc inventoryV1: %v\n", err))
		}
		closer.AddNamed("InventoryV1 grpc client", func(ctx context.Context) error {
			return conn.Close()
		})

		generatedClient := generatedInventoryV1.NewInventoryServiceClient(conn)

		d.inventoryClient = inventoryV1.NewClient(generatedClient)
	}

	return d.inventoryClient
}

func (d *diContainer) PaymentClient() grpcClient.PaymentClient {
	if d.paymentClient == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().PaymentGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("❌ не удалось подключиться к grpc paymentV1: %v\n", err))
		}
		closer.AddNamed("PaymentV1 grpc client", func(ctx context.Context) error {
			return conn.Close()
		})

		generatedClient := generatedPaymentV1.NewPaymentServiceClient(conn)

		d.paymentClient = paymentV1.NewClient(generatedClient)
	}

	return d.paymentClient
}

func (d *diContainer) MigrationRunner() migrator.Migrator {
	if d.migrationRunner == nil {
		d.migrationRunner = pg.NewMigrator(d.StdLibPgClient(), config.AppConfig().Postgres.MigrationDir())
	}

	return d.migrationRunner
}

func (d *diContainer) StdLibPgClient() *sql.DB {
	if d.stdLibPgCon == nil {
		d.stdLibPgCon = stdlib.OpenDB(*d.PgPoolConfig().ConnConfig)
	}

	return d.stdLibPgCon
}

func (d *diContainer) PgPoolConfig() *pgxpool.Config {
	if d.pgPoolCfg == nil {
		pgPoolCfg, err := pgxpool.ParseConfig(config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Sprintf("❌ не удалось разобрать конфигурацию PostgreSQL: %s\n", err.Error()))
		}

		d.pgPoolCfg = pgPoolCfg
	}

	return d.pgPoolCfg
}

func (d *diContainer) PgPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {
		pool, err := pgxpool.NewWithConfig(ctx, d.PgPoolConfig())
		if err != nil {
			panic(fmt.Sprintf("❌ не удалось создать пул PG: %v\n", err))
		}

		err = pool.Ping(ctx)
		if err != nil {
			panic(fmt.Sprintf("❌ не удалось отправить пинг клиенту PG: %s\n", err.Error()))
		}

		closer.AddNamed("PostgresDB client", func(ctx context.Context) error {
			pool.Close()
			return nil
		})

		d.pgPool = pool
	}

	return d.pgPool
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.consumerGroup.Close()
		})

		d.consumerGroup = consumerGroup
	}

	return d.consumerGroup
}

func (d *diContainer) OrderAssembledConsumer() wrapKafka.Consumer {
	if d.orderAssembledConsumer == nil {
		d.orderAssembledConsumer = wrapKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().OrderAssembledConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderAssembledConsumer
}

func (d *diContainer) OrderAssembledDecoder() kafkaConverter.OrderAssembledDecoder {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = decoder.NewOrderAssembledDecoder()
	}

	return d.orderAssembledDecoder
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) OrderPaidProducer() wrapKafka.Producer {
	if d.orderPaidProducer == nil {
		d.orderPaidProducer = wrapKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderPaidProducer.Topic(),
			logger.Logger(),
		)
	}

	return d.orderPaidProducer
}
