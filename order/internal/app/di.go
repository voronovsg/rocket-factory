package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1Impl "github.com/voronovsg/rocket-factory/order/internal/api/order/v1"
	grpcClient "github.com/voronovsg/rocket-factory/order/internal/client/grpc"
	inventoryV1 "github.com/voronovsg/rocket-factory/order/internal/client/grpc/inventory/v1"
	paymentV1 "github.com/voronovsg/rocket-factory/order/internal/client/grpc/payment/v1"
	"github.com/voronovsg/rocket-factory/order/internal/config"
	"github.com/voronovsg/rocket-factory/order/internal/repository"
	orderRepository "github.com/voronovsg/rocket-factory/order/internal/repository/order"
	"github.com/voronovsg/rocket-factory/order/internal/service"
	orderService "github.com/voronovsg/rocket-factory/order/internal/service/order"
	"github.com/voronovsg/rocket-factory/platform/pkg/closer"
	"github.com/voronovsg/rocket-factory/platform/pkg/migrator"
	"github.com/voronovsg/rocket-factory/platform/pkg/migrator/pg"
	orderV1 "github.com/voronovsg/rocket-factory/shared/pkg/openapi/order/v1"
	generatedInventoryV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/inventory/v1"
	generatedPaymentV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderV1API      orderV1.Handler
	orderService    service.OrderService
	orderRepository repository.OrderRepository

	inventoryClient grpcClient.InventoryClient
	paymentClient   grpcClient.PaymentClient
	migrationRunner migrator.Migrator

	pgPool      *pgxpool.Pool
	pgPoolCfg   *pgxpool.Config
	stdLibPgCon *sql.DB
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
		)
	}

	return d.orderService
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
