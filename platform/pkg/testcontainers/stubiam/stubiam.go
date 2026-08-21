package stubiam

import (
	"context"

	"github.com/voronovsg/rocket-factory/platform/pkg/testcontainers/app"
)

type Container struct {
	app *app.Container
	cfg *Config
}

func NewContainer(ctx context.Context, opts ...Option) (*Container, error) {
	cfg := buildConfig(opts...)

	appContainer, err := app.NewContainer(ctx,
		app.WithName(cfg.ContainerName),
		app.WithPort(cfg.Port),
		app.WithDockerfile(cfg.DockerfileDir, cfg.Dockerfile),
		app.WithNetwork(cfg.NetworkName),
		app.WithEnv(map[string]string{
			grpcHostKey: grpcHost,
			grpcPortKey: cfg.Port,
		}),
		app.WithLogOutput(cfg.LogOutput),
		app.WithStartupWait(cfg.StartupWait),
		app.WithLogger(cfg.Logger),
	)
	if err != nil {
		return nil, err
	}

	return &Container{
		app: appContainer,
		cfg: cfg,
	}, nil
}

// HostName возвращает имя контейнера для использования в docker-сети (IAM_GRPC_HOST).
func (c *Container) HostName() string {
	return c.cfg.ContainerName
}

// Port возвращает gRPC-порт stub IAM (IAM_GRPC_PORT).
func (c *Container) Port() string {
	return c.cfg.Port
}

// Address возвращает host:port для подключения с хоста (из тестового процесса).
func (c *Container) Address() string {
	return c.app.Address()
}

func (c *Container) Terminate(ctx context.Context) error {
	return c.app.Terminate(ctx)
}
