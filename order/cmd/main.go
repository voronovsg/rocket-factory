package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderAPIV1 "github.com/voronovsg/rocket-factory/order/internal/api/order/v1"
	inventoryClient "github.com/voronovsg/rocket-factory/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/voronovsg/rocket-factory/order/internal/client/grpc/payment/v1"
	"github.com/voronovsg/rocket-factory/order/internal/config"
	"github.com/voronovsg/rocket-factory/order/internal/migrator"
	orderRepo "github.com/voronovsg/rocket-factory/order/internal/repository/order"
	orderSrv "github.com/voronovsg/rocket-factory/order/internal/service/order"
	genOrderV1 "github.com/voronovsg/rocket-factory/shared/pkg/openapi/order/v1"
	genInventoryV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/inventory/v1"
	genPaymentV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/payment/v1"
)

const configPath = "deploy/compose/order/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		log.Printf("failed to load config: %v\n", err)
	}
	ctx := context.Background()

	paymentConn, err := grpc.NewClient(config.AppConfig().PaymentGRPC.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("❌ failed to connect to PaymentServer: %v\n", err)
		return
	}
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("❌ failed to close PaymentServer connection: %v", cerr)
		}
	}()

	inventoryConn, err := grpc.NewClient(config.AppConfig().InventoryGRPC.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("❌ failed to connect to InventoryServer: %v\n", err)
		return
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("❌ failed to close InventoryServer connection: %v", cerr)
		}
	}()

	genInventoryV1Client := genInventoryV1.NewInventoryServiceClient(inventoryConn)
	genPaymentV1Client := genPaymentV1.NewPaymentServiceClient(paymentConn)

	inventoryV1Client := inventoryClient.NewClient(genInventoryV1Client)
	paymentV1Client := paymentClient.NewClient(genPaymentV1Client)

	poolConf, err := pgxpool.ParseConfig(config.AppConfig().Postgres.URI())
	if err != nil {
		log.Printf("failed to parse db config: %v\n", err)
		return
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolConf)
	if err != nil {
		log.Printf("failed to create pool: %v\n", err)
		return
	}
	defer pool.Close()

	migrationRunner := migrator.New(stdlib.OpenDB(*poolConf.ConnConfig), config.AppConfig().Postgres.MigrationDir())
	err = migrationRunner.Up()
	if err != nil {
		log.Printf("failed to run migrations: %v\n", err)
		return
	}

	orderRepository := orderRepo.NewRepository(pool)
	orderService := orderSrv.NewService(orderRepository, inventoryV1Client, paymentV1Client)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	srv := orderAPIV1.NewAPI(orderService)

	handler, err := genOrderV1.NewServer(srv)
	if err != nil {
		log.Printf("failed to create ogen handler: %v\n", err)
		return
	}

	r.Mount("/", handler)

	server := &http.Server{
		ReadTimeout: config.AppConfig().OrderHTTP.ReadTimeout(),
		Addr:        config.AppConfig().OrderHTTP.Address(),
		Handler:     r,
	}

	go func() {
		log.Println("🚀 OrderService HTTP API running on :8080")
		if errServer := server.ListenAndServe(); errServer != nil && !errors.Is(errServer, http.ErrServerClosed) {
			log.Printf("server error: %v\n", errServer)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("graceful shutdown failed: %v\n", err)
		return
	}

	log.Println("✅ Server stopped")
}
