package user

import (
	"github.com/voronovsg/rocket-factory/iam/internal/repository"
	def "github.com/voronovsg/rocket-factory/iam/internal/service"
)

var _ def.UserService = (*service)(nil)

type service struct {
	userRepository repository.UserRepository
}

func NewService(userRepository repository.UserRepository) *service {
	return &service{
		userRepository: userRepository,
	}
}
