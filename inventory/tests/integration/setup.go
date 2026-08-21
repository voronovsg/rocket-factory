//go:build integration
// +build integration

package integration

import (
	"context"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
	"github.com/voronovsg/rocket-factory/platform/pkg/testcontainers"
	"github.com/voronovsg/rocket-factory/platform/pkg/testcontainers/app"
	"github.com/voronovsg/rocket-factory/platform/pkg/testcontainers/mongo"
	"github.com/voronovsg/rocket-factory/platform/pkg/testcontainers/network"
	"github.com/voronovsg/rocket-factory/platform/pkg/testcontainers/path"
	"github.com/voronovsg/rocket-factory/platform/pkg/testcontainers/stubiam"
)

const (
	inventoryAppName    = "inventory-app"
	inventoryAppPort    = "50051"
	inventoryDockerfile = "deploy/docker/inventory/Dockerfile"

	loggerLevelKey  = "LOGGER_LEVEL"
	loggerAsJsonKey = "LOGGER_AS_JSON"
	grpcHostKey     = "GRPC_HOST"
	grpcPortKey     = "GRPC_PORT"
	httpHostKey     = "HTTP_HOST"
	httpPortKey     = "HTTP_PORT"
	iamGRPCHostKey  = "IAM_GRPC_HOST"
	iamGRPCPortKey  = "IAM_GRPC_PORT"

	loggerLevelValue  = "info"
	loggerAsJsonValue = "true"
	grpcHostValue     = "0.0.0.0"
	httpHostValue     = "0.0.0.0"
	httpPortValue     = "8081"
	startupTimeout    = 3 * time.Minute
)

type TestEnvironment struct {
	Network *network.Network
	Mongo   *mongo.Container
	IAMStub *stubiam.Container
	App     *app.Container
}

// setupTestEnvironment подготавливает сеть, контейнеры и возвращает структуру с ресурсами
func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "🚀 Подготовка тестового окружения...")

	// Создаём общую Docker-сеть
	generatedNetwork, err := network.NewNetwork(ctx, projectName)
	if err != nil {
		logger.Fatal(ctx, "❌ не удалось создать общую сеть", zap.Error(err))
	}
	logger.Info(ctx, "✅ Сеть успешно создана")

	// Получаем переменные окружения для MongoDB с проверкой на наличие
	mongoUsername := getEnvWithLogging(ctx, testcontainers.MongoUsernameKey)
	mongoPassword := getEnvWithLogging(ctx, testcontainers.MongoPasswordKey)
	mongoImageName := getEnvWithLogging(ctx, testcontainers.MongoImageNameKey)

	// Запускаем контейнер с MongoDB
	generatedMongo, err := mongo.NewContainer(ctx,
		mongo.WithNetworkName(generatedNetwork.Name()),
		mongo.WithContainerName(testcontainers.MongoContainerName),
		mongo.WithImageName(mongoImageName),
		mongo.WithAuth(mongoUsername, mongoPassword),
		mongo.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "не удалось запустить контейнер MongoDB", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер MongoDB успешно запущен")

	projectRoot := path.GetProjectRoot()

	iamStubContainer, err := stubiam.NewContainer(ctx,
		stubiam.WithNetworkName(generatedNetwork.Name()),
		stubiam.WithLogOutput(os.Stdout),
		stubiam.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Mongo: generatedMongo})
		logger.Fatal(ctx, "не удалось запустить stub IAM контейнер", zap.Error(err))
	}
	logger.Info(ctx, "✅ Stub IAM контейнер успешно запущен")

	// Запускаем контейнер с приложением
	appEnv := map[string]string{
		// MongoDB переменные
		testcontainers.MongoHostKey:     generatedMongo.Config().ContainerName,
		testcontainers.MongoPortKey:     testcontainers.MongoPort,
		testcontainers.MongoDatabaseKey: generatedMongo.Config().Database,
		testcontainers.MongoUsernameKey: generatedMongo.Config().Username,
		testcontainers.MongoPasswordKey: generatedMongo.Config().Password,
		testcontainers.MongoAuthDBKey:   generatedMongo.Config().AuthDB,

		// Логгер и GRPC переменные
		loggerLevelKey:  loggerLevelValue,
		loggerAsJsonKey: loggerAsJsonValue,
		grpcHostKey:     grpcHostValue,
		grpcPortKey:     inventoryAppPort,
		httpHostKey:     httpHostValue,
		httpPortKey:     httpPortValue,
		iamGRPCHostKey:  iamStubContainer.HostName(),
		iamGRPCPortKey:  iamStubContainer.Port(),
	}

	// Создаем настраиваемую стратегию ожидания с увеличенным таймаутом
	waitStrategy := wait.ForListeningPort(inventoryAppPort + "/tcp").
		WithStartupTimeout(startupTimeout)

	appContainer, err := app.NewContainer(ctx,
		app.WithName(inventoryAppName),
		app.WithPort(inventoryAppPort),
		app.WithDockerfile(projectRoot, inventoryDockerfile),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithEnv(appEnv),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(waitStrategy),
		app.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{
			Network: generatedNetwork,
			Mongo:   generatedMongo,
			IAMStub: iamStubContainer,
		})
		logger.Fatal(ctx, "не удалось запустить контейнер приложения", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер приложения успешно запущен")

	logger.Info(ctx, "🎉 Тестовое окружение готово")
	return &TestEnvironment{
		Network: generatedNetwork,
		Mongo:   generatedMongo,
		IAMStub: iamStubContainer,
		App:     appContainer,
	}
}

// getEnvWithLogging возвращает значение переменной окружения с логированием
func getEnvWithLogging(ctx context.Context, key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Warn(ctx, "Переменная окружения не установлена", zap.String("key", key))
	}
	return value
}
