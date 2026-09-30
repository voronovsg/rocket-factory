package payment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
)

type ServiceSuite struct {
	suite.Suite
	ctx     context.Context //nolint:containedctx
	service *service
}

func (s *ServiceSuite) SetupTest() {
	logger.SetNopLogger()
	s.ctx = context.Background()
	s.service = NewService()
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
