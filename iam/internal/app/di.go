package app

import (
	"context"
	"database/sql"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	authV1Impl "github.com/voronovsg/rocket-factory/iam/internal/api/auth/v1"
	userV1Impl "github.com/voronovsg/rocket-factory/iam/internal/api/user/v1"
	"github.com/voronovsg/rocket-factory/iam/internal/config"
	"github.com/voronovsg/rocket-factory/iam/internal/repository"
	sessionRepository "github.com/voronovsg/rocket-factory/iam/internal/repository/session"
	userRepository "github.com/voronovsg/rocket-factory/iam/internal/repository/user"
	"github.com/voronovsg/rocket-factory/iam/internal/service"
	authService "github.com/voronovsg/rocket-factory/iam/internal/service/auth"
	userService "github.com/voronovsg/rocket-factory/iam/internal/service/user"
	"github.com/voronovsg/rocket-factory/platform/pkg/cache"
	"github.com/voronovsg/rocket-factory/platform/pkg/cache/redis"
	"github.com/voronovsg/rocket-factory/platform/pkg/closer"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
	"github.com/voronovsg/rocket-factory/platform/pkg/migrator"
	"github.com/voronovsg/rocket-factory/platform/pkg/migrator/pg"
	authV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/auth/v1"
	userV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/user/v1"
)

type diContainer struct {
	authV1API authV1.AuthServiceServer
	userV1API userV1.UserServiceServer

	authService     service.AuthService
	userService     service.UserService
	migrationRunner migrator.Migrator

	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
	pgPool            *pgxpool.Pool
	pgPoolCfg         *pgxpool.Config
	stdLibPgCon       *sql.DB

	redisPool   *redigo.Pool
	redisClient cache.RedisClient
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) AuthV1API(ctx context.Context) authV1.AuthServiceServer {
	if d.authV1API == nil {
		d.authV1API = authV1Impl.NewAPI(d.AuthService(ctx))
	}

	return d.authV1API
}

func (d *diContainer) UserV1API(ctx context.Context) userV1.UserServiceServer {
	if d.userV1API == nil {
		d.userV1API = userV1Impl.NewAPI(d.UserService(ctx))
	}

	return d.userV1API
}

func (d *diContainer) AuthService(ctx context.Context) service.AuthService {
	if d.authService == nil {
		d.authService = authService.NewService(
			d.UserRepository(ctx),
			d.SessionRepository(),
			config.AppConfig().Session,
		)
	}

	return d.authService
}

func (d *diContainer) UserService(ctx context.Context) service.UserService {
	if d.userService == nil {
		d.userService = userService.NewService(d.UserRepository(ctx))
	}

	return d.userService
}

func (d *diContainer) MigrationRunner() migrator.Migrator {
	if d.migrationRunner == nil {
		d.migrationRunner = pg.NewMigrator(d.StdLibPgClient(), config.AppConfig().Postgres.MigrationDirectory())
	}

	return d.migrationRunner
}

func (d *diContainer) UserRepository(ctx context.Context) repository.UserRepository {
	if d.userRepository == nil {
		d.userRepository = userRepository.NewRepository(d.PgPool(ctx))
	}

	return d.userRepository
}

func (d *diContainer) SessionRepository() repository.SessionRepository {
	if d.sessionRepository == nil {
		d.sessionRepository = sessionRepository.NewRepository(d.RedisClient())
	}

	return d.sessionRepository
}

func (d *diContainer) PgPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {
		pool, err := pgxpool.NewWithConfig(ctx, d.PgPoolConfig())
		if err != nil {
			panic(fmt.Sprintf("failed to create pg pool: %v\n", err))
		}

		err = pool.Ping(ctx)
		if err != nil {
			panic(fmt.Sprintf("failed to ping pg client: %s\n", err.Error()))
		}

		closer.AddNamed("PostgresDB client", func(ctx context.Context) error {
			pool.Close()
			return nil
		})

		d.pgPool = pool
	}

	return d.pgPool
}

func (d *diContainer) PgPoolConfig() *pgxpool.Config {
	if d.pgPoolCfg == nil {
		pgPoolCfg, err := pgxpool.ParseConfig(config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Sprintf("failed to parse postgres config: %s\n", err.Error()))
		}

		d.pgPoolCfg = pgPoolCfg
	}

	return d.pgPoolCfg
}

func (d *diContainer) StdLibPgClient() *sql.DB {
	if d.stdLibPgCon == nil {
		d.stdLibPgCon = stdlib.OpenDB(*d.PgPoolConfig().ConnConfig)
	}

	return d.stdLibPgCon
}

func (d *diContainer) RedisPool() *redigo.Pool {
	if d.redisPool == nil {
		d.redisPool = &redigo.Pool{
			MaxIdle:     config.AppConfig().Redis.MaxIdle(),
			IdleTimeout: config.AppConfig().Redis.IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", config.AppConfig().Redis.Address())
			},
		}
	}

	return d.redisPool
}

func (d *diContainer) RedisClient() cache.RedisClient {
	if d.redisClient == nil {
		d.redisClient = redis.NewClient(d.RedisPool(), logger.Logger(), config.AppConfig().Redis.ConnectionTimeout())
	}

	return d.redisClient
}
