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

// Get returns the value of the key
func (r *PostgresRepository) LoginUserPrecheck(username string) (string, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return "", err
	}
	return user.Salt, nil
}
