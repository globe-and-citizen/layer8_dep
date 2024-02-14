package interfaces

import (
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/models"
)

type IRepository interface {
	RegisterUser(req dto.RegisterUserDTO) error
	LoginPreCheckUser(req dto.LoginPrecheckDTO) (string, string, error)
	LoginPreCheckClient(req dto.LoginPrecheckDTO) (string, string, error)
	LoginUser(req dto.LoginUserDTO) (models.User, error)
	LoginClient(req dto.LoginClientDTO) (models.Client, error)
	ProfileUser(userID uint) (models.User, []models.UserMetadata, error)
	ProfileClient(userID string) (models.Client, error)
	VerifyEmail(userID uint) error
	UpdateDisplayName(userID uint, req dto.UpdateDisplayNameDTO) error
	RegisterClient(req dto.RegisterClientDTO) error
	GetClientData(clientName string) (models.Client, error)
}
