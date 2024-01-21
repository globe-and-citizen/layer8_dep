package service

import (
	"globe-and-citizen/layer8/server/resource_server/dto"
	interfaces "globe-and-citizen/layer8/server/resource_server/interfaces"

	"github.com/go-playground/validator/v10"
)

type service struct {
	repository interfaces.IRepository
}

// Newservice creates a new instance of service
func NewService(repo interfaces.IRepository) interfaces.IService {
	return &service{
		repository: repo,
	}
}

// RegisterUser registers a new user
func (s *service) RegisterUser(req dto.RegisterUserDTO) error {
	// validate request
	if err := validator.New().Struct(req); err != nil {
		return err
	}
	return s.repository.RegisterUser(req)

}
