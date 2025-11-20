package user

import (
	"testing"

	"github.com/stretchr/testify/suite"

	repositoryMocks "backend/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite
	userRepository *repositoryMocks.MockUserRepository
	service        *Service
}

func (s *ServiceSuite) SetupTest() {
	s.userRepository = repositoryMocks.NewMockUserRepository(s.T())
	s.service = NewService(s.userRepository)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
