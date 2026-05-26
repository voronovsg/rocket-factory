package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	grpcMocks "github.com/voronovsg/rocket-factory/order/internal/client/grpc/mocks"
	repoMocks "github.com/voronovsg/rocket-factory/order/internal/repository/mocks"
	serviceMocks "github.com/voronovsg/rocket-factory/order/internal/service/mocks"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
)

type ServiceSuite struct {
	suite.Suite
	mockOrderRepository      *repoMocks.MockOrderRepository
	mockInventoryClient      *grpcMocks.MockInventoryClient
	mockPaymentClient        *grpcMocks.MockPaymentClient
	mockOrderProducerService *serviceMocks.MockOrderProducerService
	service                  *service
	ctx                      context.Context //nolint:containedctx
}

func (s *ServiceSuite) SetupTest() {
	logger.SetNopLogger()
	s.mockOrderRepository = repoMocks.NewMockOrderRepository(s.T())
	s.mockInventoryClient = grpcMocks.NewMockInventoryClient(s.T())
	s.mockPaymentClient = grpcMocks.NewMockPaymentClient(s.T())
	s.mockOrderProducerService = serviceMocks.NewMockOrderProducerService(s.T())
	s.service = NewService(s.mockOrderRepository, s.mockInventoryClient, s.mockPaymentClient, s.mockOrderProducerService)
	s.ctx = context.Background()
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
