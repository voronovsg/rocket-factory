package stubiam

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
	"github.com/voronovsg/rocket-factory/platform/pkg/testcontainers"
	"github.com/voronovsg/rocket-factory/platform/pkg/testcontainers/path"
)

const (
	defaultStartupTimeout = 3 * time.Minute

	grpcHostKey = "GRPC_HOST"
	grpcPortKey = "GRPC_PORT"
	grpcHost    = "0.0.0.0"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type Config struct {
	ContainerName string
	NetworkName   string
	Port          string
	DockerfileDir string
	Dockerfile    string
	LogOutput     io.Writer
	StartupWait   wait.Strategy
	Logger        Logger
}

func buildConfig(opts ...Option) *Config {
	cfg := &Config{
		ContainerName: testcontainers.IAMStubContainerName,
		Port:          testcontainers.IAMStubPort,
		DockerfileDir: path.GetProjectRoot(),
		Dockerfile:    testcontainers.IAMStubDockerfile,
		LogOutput:     os.Stdout,
		StartupWait: wait.ForListeningPort(testcontainers.IAMStubPort + "/tcp").
			WithStartupTimeout(defaultStartupTimeout),
		Logger: &logger.NoopLogger{},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}
