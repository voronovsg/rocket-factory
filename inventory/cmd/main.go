package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	inventoryV1 "github.com/voronovsg/rocket-factory/inventory/internal/api/inventory/v1"
	"github.com/voronovsg/rocket-factory/inventory/internal/config"
	"github.com/voronovsg/rocket-factory/inventory/internal/interceptor/logger"
	"github.com/voronovsg/rocket-factory/inventory/internal/interceptor/validate"
	partRepo "github.com/voronovsg/rocket-factory/inventory/internal/repository/part"
	partSrv "github.com/voronovsg/rocket-factory/inventory/internal/service/part"
	genInventoryV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/inventory/v1"
)

const configPath = "deploy/compose/inventory/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		log.Printf("failed to load config: %v\n", err)
	}
	grpcAddr := config.AppConfig().InventoryGRPC.Address()
	httpAddr := config.AppConfig().InventoryHTTP.Address()
	ctx := context.Background()
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Printf("❌ failed to listen: %v\n", err)
		return
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
	if err != nil {
		log.Printf("❌ failed to connect to MongoDB: %v\n", err)
		return
	}
	defer func() {
		if cerr := client.Disconnect(ctx); cerr != nil {
			log.Printf("❌ failed to disconnect from MongoDB: %v\n", cerr)
		}
	}()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Printf("❌ failed to ping MongoDB: %v\n", err)
		return
	}
	log.Println("Successful connection to MongoDB")
	db := client.Database("inventory")

	repo := partRepo.NewRepository(db)
	err = repo.InitGenParts(ctx)
	if err != nil {
		log.Printf("❌ failed to init gen parts: %v\n", err)
		return
	}
	service := partSrv.NewService(repo)
	api := inventoryV1.NewAPI(service)

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			validate.UnaryValidateInterceptor(),
			logger.UnaryLoggerInterceptor(),
		),
	)
	reflection.Register(s)
	genInventoryV1.RegisterInventoryServiceServer(s, api)

	go func() {
		log.Printf("🚀 gRPC server listening on %v\n", grpcAddr)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("❌ failed to serve: %v\n", err)
			return
		}
	}()

	var gwServer *http.Server
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		err = genInventoryV1.RegisterInventoryServiceHandlerFromEndpoint(
			ctx,
			mux,
			grpcAddr,
			opts,
		)
		if err != nil {
			log.Printf("❌ failed to register gateway: %v\n", err)
			return
		}

		fileServer := http.FileServer(http.Dir("../shared/api/inventory/v1"))
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

		gwServer = &http.Server{
			Addr:              httpAddr,
			Handler:           httpMux,
			ReadHeaderTimeout: 10 * time.Second,
		}
		log.Printf("🌐 HTTP server with gRPC-Gateway and Swagger UI listening on %v\n", httpAddr)
		err = gwServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ failed to serve HTTP: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down servers...")

	if gwServer != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := gwServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("❌ HTTP server shutdown error: %v", err)
		}
		log.Println("✅ HTTP server stopped")
	}

	s.GracefulStop()
	log.Println("✅ Server stopped")
}
