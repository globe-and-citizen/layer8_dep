package repository

import (
	"globe-and-citizen/layer8/proxy/config"
	"globe-and-citizen/layer8/proxy/models"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository() *PostgresRepository {
	return &PostgresRepository{
		db: config.DB,
	}
}

func (r *PostgresRepository) LoginUserPrecheck(username string) (string, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return "", err
	}
	return user.Salt, nil
}

func (r *PostgresRepository) GetUser(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return &models.User{}, err
	}
	return &user, nil
}

func (r *PostgresRepository) GetUserByID(id int64) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return &models.User{}, err
	}
	return &user, nil
}

func (r *PostgresRepository) SetClient(client *models.Client) error {
	err := r.db.Create(&client).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) GetClient(id string) (*models.Client, error) {
	var client models.Client
	err := r.db.Where("id = ?", id).First(&client).Error
	if err != nil {
		return &models.Client{}, err
	}
	return &client, nil
}
