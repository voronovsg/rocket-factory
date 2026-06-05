package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/voronovsg/rocket-factory/inventory/internal/config"
	"github.com/voronovsg/rocket-factory/platform/pkg/closer"
	"github.com/voronovsg/rocket-factory/platform/pkg/grpc/health"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
	grpcMiddleware "github.com/voronovsg/rocket-factory/platform/pkg/middleware/grpc"
	"github.com/voronovsg/rocket-factory/platform/pkg/testcontainers/path"
	inventoryV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/inventory/v1"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	httpServer  *http.Server
	listener    net.Listener
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
		if err := a.runGRPCServer(ctx); err != nil {
			errCh <- errors.Errorf("gRPC server crashed %v", err)
		}
	}()

	go func() {
		if err := a.runHTTPServer(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- errors.Errorf("HTTP server crashed %v", err)
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
		a.initCloser,
		a.initListener,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initData,
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

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initListener(_ context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().InventoryGRPC.Address())
	if err != nil {
		return err
	}
	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		lerr := listener.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}

		return nil
	})
	a.listener = listener

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	authInterceptor := grpcMiddleware.NewAuthInterceptor(a.diContainer.GeneratedIAMClient())
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(authInterceptor.Unary()))
	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})
	reflection.Register(a.grpcServer)
	health.RegisterService(a.grpcServer)
	inventoryV1.RegisterInventoryServiceServer(a.grpcServer, a.diContainer.InventoryV1API(ctx))

	return nil
}

func (a *App) initData(ctx context.Context) error {
	return a.diContainer.PartRepository(ctx).InitGenParts(ctx)
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC InventoryService server listening on %s", config.AppConfig().InventoryGRPC.Address()))
	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	err := inventoryV1.RegisterInventoryServiceHandlerFromEndpoint(
		ctx,
		mux,
		config.AppConfig().InventoryGRPC.Address(),
		opts,
	)
	if err != nil {
		return err
	}

	staticRoot := os.Getenv("INVENTORY_HTTP_STATIC_DIR")
	if staticRoot == "" {
		staticRoot = filepath.Join(path.GetProjectRoot(), "shared", "api", "inventory", "v1")
	}
	fileServer := http.FileServer(http.Dir(staticRoot))
	httpMux := http.NewServeMux()
	httpMux.Handle("/api/", mux)
	httpMux.Handle("/swagger-ui.html", fileServer)
	httpMux.Handle("/inventory.swagger.json", fileServer)
	httpMux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/swagger-ui.html", http.StatusMovedPermanently)
			return
		}
		fileServer.ServeHTTP(w, r)
	}))

	a.httpServer = &http.Server{
		Addr:              config.AppConfig().InventoryHTTP.Address(),
		Handler:           httpMux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		if err := a.httpServer.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	})

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🌐 HTTP server with gRPC-Gateway and Swagger UI listening on %v\n", config.AppConfig().InventoryHTTP.Address()))

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
