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

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/voronovsg/rocket-factory/payment/internal/interceptor/validate"
	paymentV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/payment/v1"
)

const (
	grpcAddr = "localhost:50052"
	httpAddr = "localhost:8082"
)

type PaymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

// PayOrder обрабатывает команду на оплату и возвращает UUID транзакции
func (s *PaymentService) PayOrder(_ context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v\n", err)
	}
	log.Printf("Order UUID: %s\nUser UUID: %s\nPayment Method: %s",
		req.GetOrderUuid(),
		req.GetUserUuid(),
		req.GetPaymentMethod().Enum().String())

	transactionUUID := uuid.NewString()
	log.Printf("Оплата прошла успешно, Transaction UUID: %s", transactionUUID)

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(validate.UnaryValidateInterceptor()),
	)
	reflection.Register(s)

	service := &PaymentService{}
	paymentV1.RegisterPaymentServiceServer(s, service)

	go func() {
		log.Printf("🚀 gRPC PaymentService server listening on %s\n", grpcAddr)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	var gwServer *http.Server
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		err = paymentV1.RegisterPaymentServiceHandlerFromEndpoint(
			ctx,
			mux,
			grpcAddr,
			opts,
		)
		if err != nil {
			log.Printf("Failed to register gateway: %v\n", err)
			return
		}

		fileServer := http.FileServer(http.Dir("../shared/api/payment/v1"))
		httpMux := http.NewServeMux()
		httpMux.Handle("/api/", mux)
		httpMux.Handle("/swagger-ui.html", fileServer)
		httpMux.Handle("/payment.swagger.json", fileServer)
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
			log.Printf("Failed to serve HTTP: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")

	if gwServer != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := gwServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
		log.Println("✅ HTTP server stopped")
	}

	s.GracefulStop()
	log.Println("✅ Server stopped")
}
