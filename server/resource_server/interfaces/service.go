package interfaces

import "globe-and-citizen/layer8/server/resource_server/dto"

type IService interface {
	RegisterUser(req dto.RegisterUserDTO) error
}
