package repository

import (
	"fmt"
	"globe-and-citizen/layer8/server/models"
	"time"

	uuid "github.com/google/uuid"
)

var TheInMemoryRepository *InMemoryRepository

func init() {
	// This function is only run once no matter how many times the pkg is imported
	TheInMemoryRepository = NewInMemoryRepository()
	TheInMemoryRepository.Clients["1"] = &models.Client{
		ID:          "abc123isNotAnId",
		Secret:      "7af4a43a0439ea61f044f6bbc1e88af21f66126a6cad8d5b08777a98d7b3457a",
		Name:        "MyMemoryClient",
		RedirectURI: "http://localhost:5173/oauth2/callback",
	}
	TheInMemoryRepository.Users["1"] = &models.User{
		ID:        uint(1),
		Email:     "tester_chester@gmail.com",
		Username:  "tester",
		Password:  "12341234",
		FirstName: "Tester",
		LastName:  "Chester",
		Salt:      "64e7f34a24237c0f36db144aed97920c542184f1c129185444f9596a5252a4f1",
	}
	TheInMemoryRepository.UserMetadata["1"] = &models.UserMetadata{
		ID:     int64(1),
		UserID: int64(1),
		Key:    "email_verified",
		Value:  "false",
	}

	TheInMemoryRepository.UserMetadata["2"] = &models.UserMetadata{
		ID:     int64(2),
		UserID: int64(1),
		Key:    "country",
		Value:  "Canada",
	}

	TheInMemoryRepository.UserMetadata["3"] = &models.UserMetadata{
		ID:     int64(3),
		UserID: int64(1),
		Key:    "display_name",
		Value:  "t_chester",
	}

}

type InMemoryRepository struct {
	Clients      map[string]*models.Client
	Users        map[string]*models.User
	UserMetadata map[string]*models.UserMetadata
	storage      map[string][]byte
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		Clients:      make(map[string]*models.Client),
		Users:        make(map[string]*models.User),
		UserMetadata: make(map[string]*models.UserMetadata),
		storage:      make(map[string][]byte),
	}
}

func (imr *InMemoryRepository) LoginUserPrecheck(username string) (string, error) {
	fmt.Println("LoginUserPrecheck", username)
	for _, user := range imr.Users {
		if user.Username == username {
			return user.Salt, nil
		}
	}
	return "", fmt.Errorf("user not found")
}

func (imr *InMemoryRepository) GetUser(username string) (*models.User, error) {
	fmt.Println("GetUser", username)
	//var val *models.User
	for _, user := range imr.Users {
		if user.Username == username {
			return &models.User{}, nil
		}
	}
	return &models.User{}, fmt.Errorf("user not found")
}

func (imr *InMemoryRepository) GetUserByID(id int64) (*models.User, error) {
	fmt.Println("GetUserByID", id)
	// var val *models.User
	for _, user := range imr.Users {
		if int64(user.ID) == id {
			return &models.User{}, nil
		}
	}
	return &models.User{}, fmt.Errorf("user not found")
}

func (imr *InMemoryRepository) GetUserMetadata(userID int64, key string) (*models.UserMetadata, error) {
	fmt.Println("GetUserMetadata", userID)
	for _, userMetadata := range imr.UserMetadata {
		if userMetadata.ID == userID {
			if userMetadata.Key == key {
				return &models.UserMetadata{}, nil
			}
		}
	}

	return &models.UserMetadata{}, nil
}

func (imr *InMemoryRepository) SetClient(client *models.Client) error {
	fmt.Println("SetClient: ", client)
	// Check if client already exists
	for _, imrClient := range imr.Clients {
		if imrClient.Name == client.Name {
			return fmt.Errorf("client already existed")
		}
	}

	uuid := uuid.New()
	UUID := uuid.String()
	imr.Clients[UUID] = client

	return nil
}

func (imr *InMemoryRepository) GetClient(id string) (*models.Client, error) {
	fmt.Println("GetClient ", id)
	for _, client := range imr.Clients {
		if client.ID == id {
			return client, nil
		}
	}
	return &models.Client{}, fmt.Errorf("no matching client found")
}

func (imr *InMemoryRepository) SetTTL(key string, value []byte, ttl time.Duration) error {
	fmt.Println("SetTTL", key, value)
	imr.storage[key] = value
	go func() {
		time.Sleep(ttl)
		delete(imr.storage, key)
	}()
	return nil
}

func (imr *InMemoryRepository) GetTTL(key string) ([]byte, error) {
	fmt.Println("GetTTL: ", key)
	if imr.storage[key] == nil {
		return nil, fmt.Errorf("func 'GetTTL' failed. Specified key not found")
	}
	return make([]byte, 32), nil
}
