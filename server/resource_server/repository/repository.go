package repository

import (
	"globe-and-citizen/layer8/server/config"
	"globe-and-citizen/layer8/server/resource_server/dto"
	interfaces "globe-and-citizen/layer8/server/resource_server/interfaces"
	"globe-and-citizen/layer8/server/resource_server/models"
	"globe-and-citizen/layer8/server/resource_server/utils"

	"gorm.io/gorm"
)

type Repository struct {
	connection *gorm.DB
}

func NewRepository(db *gorm.DB) interfaces.IRepository {
	return &Repository{
		connection: db,
	}
}

func (r *Repository) RegisterUser(req dto.RegisterUserDTO) error {

	rmSalt := utils.GenerateRandomSalt(utils.SaltSize)
	HashedAndSaltedPass := utils.SaltAndHashPassword(req.Password, rmSalt)

	user := models.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  HashedAndSaltedPass,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Salt:      rmSalt,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return err
	}

	userMetadata := []models.UserMetadata{
		{
			UserID: user.ID,
			Key:    "email_verified",
			Value:  "false",
		},
		{
			UserID: user.ID,
			Key:    "country",
			Value:  req.Country,
		},
		{
			UserID: user.ID,
			Key:    "display_name",
			Value:  req.DisplayName,
		},
	}

	if err := config.DB.Create(&userMetadata).Error; err != nil {
		config.DB.Delete(&user)
		return err
	}

	return nil
}

func (r *Repository) RegisterClient(req dto.RegisterClientDTO) error {

	clientUUID := utils.GenerateUUID()
	clientSecret := utils.GenerateSecret(utils.SecretSize)

	rmSalt := utils.GenerateRandomSalt(utils.SaltSize)
	HashedAndSaltedPass := utils.SaltAndHashPassword(req.Password, rmSalt)

	client := models.Client{
		ID:          clientUUID,
		Secret:      clientSecret,
		Name:        req.Name,
		RedirectURI: req.RedirectURI,
		Username:  	 req.Username,
		Password:    HashedAndSaltedPass,
		Salt: 		 rmSalt,
	}

	if err := config.DB.Create(&client).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetClientData(clientName string) (models.Client, error) {
	var client models.Client
	if err := config.DB.Where("name = ?", clientName).First(&client).Error; err != nil {
		return models.Client{}, err
	}
	return client, nil
}

func (r *Repository) LoginPreCheckUser(req dto.LoginPrecheckDTO) (string, string, error) {
	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return "", "", err
	}
	return user.Username, user.Salt, nil
}

func (r *Repository) LoginPreCheckClient(req dto.LoginPrecheckDTO) (string, string, error) {
	var client models.Client
	if err := config.DB.Where("username = ?", req.Username).First(&client).Error; err != nil {
		return "", "", err
	}
	return client.Username, client.Salt, nil
}

func (r *Repository) LoginUser(req dto.LoginUserDTO) (models.User, error) {
	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *Repository) LoginClient(req dto.LoginClientDTO) (models.Client, error) {
	var client models.Client
	if err := config.DB.Where("username = ?", req.Username).First(&client).Error; err != nil {
		return models.Client{}, err
	}
	return client, nil
}

func (r *Repository) ProfileUser(userID uint) (models.User, []models.UserMetadata, error) {
	var user models.User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return models.User{}, []models.UserMetadata{}, err
	}
	var userMetadata []models.UserMetadata
	if err := config.DB.Where("user_id = ?", userID).Find(&userMetadata).Error; err != nil {
		return models.User{}, []models.UserMetadata{}, err
	}
	return user, userMetadata, nil
}

func (r *Repository) ProfileClient(userID string) (models.Client, error) {
	var client models.Client
	if err := config.DB.Where("id = ?", userID).First(&client).Error; err != nil {
		return models.Client{}, err
	}
	return client, nil
}

func (r *Repository) VerifyEmail(userID uint) error {
	return config.DB.Model(&models.UserMetadata{}).Where("user_id = ? AND key = ?", userID, "email_verified").Update("value", "true").Error
}

func (r *Repository) UpdateDisplayName(userID uint, req dto.UpdateDisplayNameDTO) error {
	return config.DB.Model(&models.UserMetadata{}).Where("user_id = ? AND key = ?", userID, "display_name").Update("value", req.DisplayName).Error
}
