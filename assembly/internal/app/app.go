package app

import (
	"context"

	"github.com/voronovsg/rocket-factory/assembly/internal/config"
	"github.com/voronovsg/rocket-factory/assembly/internal/metrics"
	"github.com/voronovsg/rocket-factory/platform/pkg/closer"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
)

type App struct {
	diContainer *diContainer
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	logger.Info(ctx, "🚀 AssemblyService running")

	return a.diContainer.OrderConsumerService().RunConsumer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initMetrics,
		a.initCloser,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	return logger.Init(
		ctx,
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
		config.AppConfig().Logger.EnableOTLP(),
		config.AppConfig().Logger.CollectorEndpoint(),
		config.AppConfig().Logger.ServiceName(),
		config.AppConfig().Logger.ServiceEnvironment(),
	)
}

func (a *App) initMetrics(ctx context.Context) error {
	err := metrics.InitMetrics(ctx, config.AppConfig().MetricServer)
	if err != nil {
		return err
	}
	closer.AddNamed("metrics", metrics.ShutdownMetrics)
	return nil
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}
