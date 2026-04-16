package part

import (
	"testing"

	"github.com/stretchr/testify/suite"

	repoMocks "github.com/voronovsg/rocket-factory/inventory/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite                                      // Встраивание testify/suite даёт все методы testify
	mockPartRepository *repoMocks.MockPartRepository // Мок репозитория - имитация работы с БД
	service            *service                      // Тестируемый объект, который использует мок
}

func (s *ServiceSuite) SetupTest() {
	s.mockPartRepository = repoMocks.NewMockPartRepository(s.T())
	s.service = NewService(s.mockPartRepository)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
