package repository

import (
	"globe-and-citizen/layer8/server/config"
	"globe-and-citizen/layer8/server/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db      *gorm.DB
	storage map[string][]byte
}

func NewPostgresRepository() *PostgresRepository {
	return &PostgresRepository{
		db:      config.DB,
		storage: make(map[string][]byte),
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

func (r *PostgresRepository) GetUserMetadata(userID int64, key string) (*models.UserMetadata, error) {
	var userMetadata models.UserMetadata
	err := r.db.Where("user_id = ? AND key = ?", userID, key).First(&userMetadata).Error
	if err != nil {
		return &models.UserMetadata{}, err
	}
	return &userMetadata, nil
}

func (r *PostgresRepository) SetClient(client *models.Client) error {
	// Check if client already exists
	var existingClient models.Client
	err := r.db.Where("id = ?", client.ID).First(&existingClient).Error
	if err == nil {
		return nil
	}
	err = r.db.Create(&client).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) GetClient(id string) (*models.Client, error) {
	if func() bool {
		var prefix string = "client:"
		return len(id) >= len(prefix) && id[0:len(prefix)] == prefix
	}() {
		id = strings.TrimPrefix(id, "client:")
	}
	var client models.Client
	err := r.db.Where("id = ?", id).First(&client).Error
	if err != nil {
		return &models.Client{}, err
	}
	return &client, nil
}

// SetTTL sets the key to hold the value for a limited time
func (r *PostgresRepository) SetTTL(key string, value []byte, ttl time.Duration) error {
	r.storage[key] = value
	go func() {
		time.Sleep(ttl)
		delete(r.storage, key)
	}()
	return nil
}

// GetTTL gets the value for the given key which has a short TTL.
func (r *PostgresRepository) GetTTL(key string) ([]byte, error) {
	return r.storage[key], nil
}
