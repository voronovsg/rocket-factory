package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1 "github.com/voronovsg/rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/payment/v1"
)

const (
	httpPort              = ":8080"
	httpServerReadTimeout = 5 * time.Second

	inventoryServiceAddr = "localhost:50051"
	paymentServiceAddr   = "localhost:50052"
)

const (
	orderStatusPendingPayment = "PENDING_PAYMENT"
	orderStatusPaid           = "PAID"
	orderStatusCancelled      = "CANCELLED"
)

type orderServer struct {
	mu     sync.RWMutex
	orders map[string]*orderData

	inventoryV1Client inventoryV1.InventoryServiceClient
	paymentV1Client   paymentV1.PaymentServiceClient
}

type orderData struct {
	UUID            string
	UserUUID        string
	PartUuids       []string
	TotalPrice      float64
	TransactionUUID *string
	PaymentMethod   *string
	Status          string
	CreatedAt       time.Time
}

func NewOrderServer(
	inventoryV1Client inventoryV1.InventoryServiceClient,
	paymentV1Client paymentV1.PaymentServiceClient,
) *orderServer {
	return &orderServer{
		orders: make(map[string]*orderData),

		inventoryV1Client: inventoryV1Client,
		paymentV1Client:   paymentV1Client,
	}
}

// CreateOrder создаёт новый заказ на основе выбранных пользователем деталей
// Получает детали из InventoryService ListParts
// Проверяет, что все детали существуют. Если хотя бы одной нет — возвращает ошибку
// Сохраняет заказ со статусом PENDING_PAYMENT
func (s *orderServer) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	partUuids := make([]string, 0, len(req.PartUuids))
	for _, partUuid := range req.PartUuids {
		partUuids = append(partUuids, partUuid.String())
	}

	res, err := s.inventoryV1Client.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: partUuids,
		},
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return &orderV1.ServiceUnavailableError{
				Message: "inventory service timeout",
			}, nil
		}

		return &orderV1.InternalServerError{
			Message: fmt.Sprintf("inventory service error: %s", err.Error()),
		}, nil
	}

	if len(res.Parts) != len(req.PartUuids) {
		return &orderV1.BadRequestError{
			Message: "some parts not found",
		}, nil
	}

	var totalPrice float64
	for _, part := range res.Parts {
		totalPrice += part.Price
	}

	newOrderUUID := uuid.NewString()

	order := &orderData{
		UUID:       newOrderUUID,
		UserUUID:   req.UserUUID.String(),
		PartUuids:  partUuids,
		TotalPrice: totalPrice,
		Status:     orderStatusPendingPayment,
		CreatedAt:  time.Now(),
	}

	s.mu.Lock()
	s.orders[newOrderUUID] = order
	s.mu.Unlock()

	log.Printf(`[Order Created]
	Order UUID: %s
	User UUID: %s
	Part UUIDs: %v
	Total Price: %f
	Status: %s
	CreatedAt: %v`, order.UUID, order.UserUUID, order.PartUuids, order.TotalPrice, order.Status, order.CreatedAt)

	return &orderV1.CreateOrderResponse{
		Order: orderV1.OrderDto{
			UUID:       uuid.MustParse(newOrderUUID),
			UserUUID:   req.UserUUID,
			PartUuids:  req.PartUuids,
			TotalPrice: order.TotalPrice,
			Status:     orderV1.OrderStatus(order.Status),
			CreatedAt:  order.CreatedAt,
		},
	}, nil
}

// PayOrder обрабатывает запрос на оплату заказа
// Проверяет существование заказа и его статус
// вызывает платёжный сервис и обновляет статус заказа при успешной оплате
func (s *orderServer) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[params.OrderUUID]
	if !ok {
		return &orderV1.NotFoundError{
			Message: "order not found",
		}, nil
	}

	if order.Status != orderStatusPendingPayment {
		return &orderV1.ConflictError{
			Message: "order cannot be paid",
		}, nil
	}

	res, err := s.paymentV1Client.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     order.UUID,
		UserUuid:      order.UserUUID,
		PaymentMethod: paymentV1.PaymentMethod(paymentV1.PaymentMethod_value[string(req.PaymentMethod)]),
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return &orderV1.ServiceUnavailableError{
				Message: "payment service timeout",
			}, nil
		}

		return &orderV1.InternalServerError{
			Message: fmt.Sprintf("payment service error: %s", err.Error()),
		}, nil
	}

	paymentMethod := string(req.PaymentMethod)

	order.TransactionUUID = &res.TransactionUuid
	order.PaymentMethod = &paymentMethod
	order.Status = orderStatusPaid

	log.Printf(`[Order Paid]
	Order UUID: %s
	User UUID: %s
	Transaction UUID: %s
	Payment Method: %s
	Status: %s`, order.UUID, order.UserUUID, *order.TransactionUUID, *order.PaymentMethod, order.Status)

	return &orderV1.PayOrderResponse{
		TransactionUUID: uuid.MustParse(res.TransactionUuid),
	}, nil
}

// GetOrderByUUID возвращает информацию о заказе по его UUID, если не найден - 404 Not Found
func (s *orderServer) GetOrderByUUID(_ context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[params.OrderUUID]
	if !ok {
		return &orderV1.NotFoundError{
			Message: "order not found",
		}, nil
	}

	partUuids := make([]uuid.UUID, 0, len(order.PartUuids))
	for _, partUuid := range order.PartUuids {
		partUuids = append(partUuids, uuid.MustParse(partUuid))
	}

	var transactionUUID orderV1.OptNilUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderV1.NewOptNilUUID(uuid.MustParse(*order.TransactionUUID))
	}

	var paymentMethod orderV1.OptPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderV1.NewOptPaymentMethod(orderV1.PaymentMethod(*order.PaymentMethod))
	}

	return &orderV1.GetOrderResponse{
		Order: orderV1.OrderDto{
			UUID:            uuid.MustParse(order.UUID),
			UserUUID:        uuid.MustParse(order.UserUUID),
			PartUuids:       partUuids,
			TotalPrice:      order.TotalPrice,
			TransactionUUID: transactionUUID,
			PaymentMethod:   paymentMethod,
			Status:          orderV1.OrderStatus(order.Status),
			CreatedAt:       order.CreatedAt,
		},
	}, nil
}

// CancelOrderByUUID отменяет заказ, если PENDING_PAYMENT — меняет статус на CANCELLED
// если PAID или CANCELLED — возвращает ошибку 409
func (s *orderServer) CancelOrderByUUID(_ context.Context, params orderV1.CancelOrderByUUIDParams) (orderV1.CancelOrderByUUIDRes, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[params.OrderUUID]
	if !ok {
		return &orderV1.NotFoundError{
			Message: "order not found",
		}, nil
	}

	if order.Status != orderStatusPendingPayment {
		return &orderV1.ConflictError{
			Message: "order cannot be canceled",
		}, nil
	}

	order.Status = orderStatusCancelled

	log.Printf(`[Order Canceled]
	Order UUID: %s
	User UUID: %s
	Status: %s`, order.UUID, order.UserUUID, order.Status)

	return &orderV1.CancelOrderByUUIDNoContent{}, nil
}

// NewError реализует интерфейс обработчика ошибок для ogen
// возвращает клиенту 500 Internal Server Error в случае непредвиденных ошибок
func (s *orderServer) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Message: err.Error(),
		},
	}
}

func main() {
	paymentConn, err := grpc.NewClient(paymentServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("❌ failed to connect to PaymentServer: %v\n", err)
		return
	}
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("❌ failed to close PaymentServer connection: %v", cerr)
		}
	}()

	inventoryConn, err := grpc.NewClient(inventoryServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("❌ failed to connect to InventoryServer: %v\n", err)
		return
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("❌ failed to close InventoryServer connection: %v", cerr)
		}
	}()

	paymentV1Client := paymentV1.NewPaymentServiceClient(paymentConn)
	inventoryV1Client := inventoryV1.NewInventoryServiceClient(inventoryConn)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	srv := NewOrderServer(inventoryV1Client, paymentV1Client)

	handler, err := orderV1.NewServer(srv)
	if err != nil {
		log.Printf("failed to create ogen handler: %v\n", err)
		return
	}

	r.Mount("/", handler)

	server := &http.Server{
		ReadTimeout: httpServerReadTimeout,
		Addr:        httpPort,
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
