package order

import (
	"github.com/voronovsg/rocket-factory/order/internal/client/grpc"
	"github.com/voronovsg/rocket-factory/order/internal/repository"
	def "github.com/voronovsg/rocket-factory/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepository repository.OrderRepository

	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient

	orderProducerService def.OrderProducerService
}

func NewService(
	orderRepository repository.OrderRepository,

	inventoryClient grpc.InventoryClient,
	paymentClient grpc.PaymentClient,

	orderProducerService def.OrderProducerService,
) *service {
	return &service{
		orderRepository: orderRepository,

		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,

		orderProducerService: orderProducerService,
	}
}
