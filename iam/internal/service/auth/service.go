package auth

import (
	"github.com/voronovsg/rocket-factory/iam/internal/config"
	"github.com/voronovsg/rocket-factory/iam/internal/repository"
	def "github.com/voronovsg/rocket-factory/iam/internal/service"
)

var _ def.AuthService = (*service)(nil)

type service struct {
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
	sessionConfig     config.SessionConfig
}

func NewService(
	userRepository repository.UserRepository,
	sessionRepository repository.SessionRepository,
	sessionConfig config.SessionConfig,
) *service {
	return &service{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		sessionConfig:     sessionConfig,
	}
}
