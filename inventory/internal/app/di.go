package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryV1Impl "github.com/voronovsg/rocket-factory/inventory/internal/api/inventory/v1"
	grpcClient "github.com/voronovsg/rocket-factory/inventory/internal/client/grpc"
	iamV1 "github.com/voronovsg/rocket-factory/inventory/internal/client/grpc/iam/v1"
	"github.com/voronovsg/rocket-factory/inventory/internal/config"
	"github.com/voronovsg/rocket-factory/inventory/internal/repository"
	partRepository "github.com/voronovsg/rocket-factory/inventory/internal/repository/part"
	"github.com/voronovsg/rocket-factory/inventory/internal/service"
	partService "github.com/voronovsg/rocket-factory/inventory/internal/service/part"
	"github.com/voronovsg/rocket-factory/platform/pkg/closer"
	authV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/auth/v1"
	inventoryV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryV1API     inventoryV1.InventoryServiceServer
	partService        service.PartService
	iamClient          grpcClient.IAMClient
	generatedIAMClient authV1.AuthServiceClient

	partRepository repository.PartRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = inventoryV1Impl.NewAPI(d.PartService(ctx))
	}

	return d.inventoryV1API
}

func (d *diContainer) PartService(ctx context.Context) service.PartService {
	if d.partService == nil {
		d.partService = partService.NewService(d.PartRepository(ctx))
	}

	return d.partService
}

func (d *diContainer) PartRepository(ctx context.Context) repository.PartRepository {
	if d.partRepository == nil {
		d.partRepository = partRepository.NewRepository(d.MongoDBHandle(ctx))
	}

	return d.partRepository
}

func (d *diContainer) IAMClient() grpcClient.IAMClient {
	if d.iamClient == nil {
		d.iamClient = iamV1.NewClient(d.GeneratedIAMClient())
	}

	return d.iamClient
}

func (d *diContainer) GeneratedIAMClient() authV1.AuthServiceClient {
	if d.generatedIAMClient == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().IAMGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to generated IAM grpc connect: %v\n", err))
		}
		closer.AddNamed("generated IAM grpc client", func(ctx context.Context) error {
			return conn.Close()
		})

		d.generatedIAMClient = authV1.NewAuthServiceClient(conn)
	}

	return d.generatedIAMClient
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("❌ failed to connect to MongoDB: %s\n", err.Error()))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("❌ failed to ping MongoDB: %v\n", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DBName())
	}

	return d.mongoDBHandle
}
