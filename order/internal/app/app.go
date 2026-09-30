package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-faster/errors"
	"go.uber.org/zap"

	"github.com/voronovsg/rocket-factory/order/internal/api/health"
	"github.com/voronovsg/rocket-factory/order/internal/config"
	"github.com/voronovsg/rocket-factory/order/internal/metrics"
	metricsMiddleware "github.com/voronovsg/rocket-factory/order/internal/middleware"
	"github.com/voronovsg/rocket-factory/platform/pkg/closer"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
	httpMiddleware "github.com/voronovsg/rocket-factory/platform/pkg/middleware/http"
	"github.com/voronovsg/rocket-factory/platform/pkg/tracing"
	generatedOrderV1 "github.com/voronovsg/rocket-factory/shared/pkg/openapi/order/v1"
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
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
	errCh := make(chan error, 2)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := a.runConsumer(ctx); err != nil {
			errCh <- errors.Errorf("consumer crashed: %v", err)
		}
	}()

	go func() {
		if err := a.runHTTPServer(ctx); err != nil {
			errCh <- errors.Errorf("http server crashed: %v", err)
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "Shutdown signal received")
	case err := <-errCh:
		logger.Error(ctx, "Component crashed, shutting down", zap.Error(err))
		cancel()
		<-ctx.Done()
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initTracer,
		a.initMetrics,
		a.initCloser,
		a.initHTTPServer,
		a.initMigration,
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

func (a *App) initTracer(ctx context.Context) error {
	err := tracing.InitTracer(ctx, config.AppConfig().Trace)
	if err != nil {
		return err
	}
	closer.AddNamed("tracer", tracing.ShutdownTracer)
	return nil
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

func (a *App) initHTTPServer(ctx context.Context) error {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Публичные маршруты без метрик и трейсинга
	r.Group(func(r chi.Router) {
		r.Get("/health", health.Handler)
	})

	handler, err := generatedOrderV1.NewServer(a.diContainer.OrderV1API(ctx))
	if err != nil {
		return err
	}

	// Защищённые маршруты с метриками, трейсингом, авторизацией
	r.Group(func(r chi.Router) {
		r.Use(metricsMiddleware.MetricsMiddleware)
		r.Use(middleware.Timeout(5 * time.Second))
		r.Use(tracing.HTTPHandlerMiddleware(config.AppConfig().Trace.ServiceName()))
		r.Use(httpMiddleware.NewAuthMiddleware(a.diContainer.GeneratedIAMClient()).Handle)

		r.Mount("/", handler)
	})

	a.httpServer = &http.Server{
		ReadTimeout:       config.AppConfig().OrderHTTP.ReadTimeout(),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		Addr:              config.AppConfig().OrderHTTP.Address(),
		Handler:           r,
	}

	closer.AddNamed("order HTTP server", func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})

	return nil
}

func (a *App) initMigration(_ context.Context) error {
	return a.diContainer.MigrationRunner().Up()
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 OrderService HTTP API launched on %s", config.AppConfig().OrderHTTP.Address()))

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 OrderService Kafka consumer running")

	err := a.diContainer.OrderConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}
