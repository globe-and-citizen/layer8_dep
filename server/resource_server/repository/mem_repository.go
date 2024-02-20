package repository

import (
	"fmt"
	serverModels "globe-and-citizen/layer8/server/models"
	"globe-and-citizen/layer8/server/resource_server/dto"
	interfaces "globe-and-citizen/layer8/server/resource_server/interfaces"
	"globe-and-citizen/layer8/server/resource_server/models"
	"globe-and-citizen/layer8/server/resource_server/utils"
	"strconv"
	"strings"
	"time"
)

type MemoryRepository struct {
	storage     map[string]map[string]string
	byteStorage map[string][]byte
}

func NewMemoryRepository() interfaces.IRepository {
	return &MemoryRepository{
		storage:     make(map[string]map[string]string),
		byteStorage: make(map[string][]byte),
	}
}

func (r *MemoryRepository) RegisterUser(req dto.RegisterUserDTO) error {
	rmSalt := utils.GenerateRandomSalt(utils.SaltSize)
	HashedAndSaltedPass := utils.SaltAndHashPassword(req.Password, rmSalt) // what if two user's use the same password?
	userID := fmt.Sprintf("%d", len(r.storage))
	r.storage[req.Password] = map[string]string{
		"user_id":        userID,
		"email":          req.Email,
		"username":       req.Username,
		"password":       HashedAndSaltedPass,
		"first_name":     req.FirstName,
		"last_name":      req.LastName,
		"country":        req.Country,
		"display_name":   req.DisplayName,
		"email_verified": "false",
	}
	r.storage[req.Username] = map[string]string{
		"salt":     rmSalt,
		"password": req.Password,
	}
	r.storage[userID] = map[string]string{
		"password": req.Password,
	}
	return nil
}

// RAVI
func (r *MemoryRepository) RegisterClient(req dto.RegisterClientDTO) error {
	if _, ok := r.storage[req.Username]; ok {
		return fmt.Errorf("Client username already registered.")
	}
	clientUUID := utils.GenerateUUID()
	clientSecret := utils.GenerateSecret(utils.SecretSize)
	r.storage[req.Username] = map[string]string{
		"id":           clientUUID,
		"secret":       clientSecret,
		"redirect_uri": req.RedirectURI,
		"username":     req.Username,
		"password":     req.Password,
	}
	return nil
}

func (r *MemoryRepository) GetClientData(clientName string) (models.Client, error) {
	if _, ok := r.storage[clientName]; !ok {
		return models.Client{}, fmt.Errorf("client not found")
	}
	client := models.Client{
		ID:          r.storage[clientName]["id"],
		Secret:      r.storage[clientName]["secret"],
		Name:        clientName,
		RedirectURI: r.storage[clientName]["redirect_uri"],
	}
	return client, nil
}

// RAVI -- Do I Need a Login client?

func (r *MemoryRepository) LoginPreCheckUser(req dto.LoginPrecheckDTO) (string, string, error) {
	if _, ok := r.storage[req.Username]["salt"]; !ok {
		return "", "", fmt.Errorf("salt not found for specified user")
	}
	return req.Username, r.storage[req.Username]["salt"], nil
}

func (r *MemoryRepository) LoginUser(req dto.LoginUserDTO) (models.User, error) {
	if _, ok := r.storage[req.Password]; !ok {
		return models.User{}, fmt.Errorf("user not found")
	}
	if r.storage[req.Password]["username"] != req.Username {
		return models.User{}, fmt.Errorf("invalid username")
	}
	UserID := r.storage[req.Password]["user_id"]
	userIdUint, err := strconv.ParseUint(UserID, 10, 32)
	if err != nil {
		return models.User{}, err
	}
	user := models.User{
		ID:        uint(userIdUint),
		Email:     r.storage[req.Password]["email"],
		Username:  r.storage[req.Password]["username"],
		Password:  r.storage[req.Password]["password"],
		FirstName: r.storage[req.Password]["first_name"],
		LastName:  r.storage[req.Password]["last_name"],
		Salt:      r.storage[req.Username]["salt"],
	}
	return user, nil
}

// RAVI
func (r *MemoryRepository) LoginClient(req dto.LoginClientDTO) (models.Client, error) {
	if _, ok := r.storage[req.Username]; !ok {
		return models.Client{}, fmt.Errorf("user not found")
	}

	if r.storage[req.Username]["password"] != req.Password {
		return models.Client{}, fmt.Errorf("invalid password")
	}

	client := models.Client{
		ID:          r.storage[req.Username]["id"],
		Secret:      r.storage[req.Username]["secret"],
		RedirectURI: r.storage[req.Username]["redirect_uri"],
		Username:    r.storage[req.Username]["username"],
	}

	fmt.Println("client: ", client)
	return client, nil
}

func (r *MemoryRepository) ProfileUser(userID uint) (models.User, []models.UserMetadata, error) {
	if _, ok := r.storage[fmt.Sprintf("%d", userID)]; !ok {
		return models.User{}, []models.UserMetadata{}, fmt.Errorf("user not found")
	}
	password := r.storage[fmt.Sprintf("%d", userID)]["password"]
	user := models.User{
		ID:        userID,
		Email:     r.storage[password]["email"],
		Username:  r.storage[password]["username"],
		Password:  r.storage[password]["password"],
		FirstName: r.storage[password]["first_name"],
		LastName:  r.storage[password]["last_name"],
		Salt:      r.storage[password]["salt"],
	}
	userMetadata := []models.UserMetadata{
		{
			Key:   "display_name",
			Value: r.storage[password]["display_name"],
		},
		{
			Key:   "country",
			Value: r.storage[password]["country"],
		},
		{
			Key:   "email_verified",
			Value: r.storage[password]["email_verified"],
		},
	}
	return user, userMetadata, nil
}

func (r *MemoryRepository) ProfileClient(username string) (models.Client, error) {
	if _, ok := r.storage[username]; !ok {
		return models.Client{}, fmt.Errorf("user not found")
	}

	client := models.Client{
		ID:          r.storage[username]["id"],
		Secret:      r.storage[username]["secret"],
		RedirectURI: r.storage[username]["redirect_uri"],
		Username:    username,
	}
	return client, nil
}

func (r *MemoryRepository) VerifyEmail(userID uint) error {
	if _, ok := r.storage[fmt.Sprintf("%d", userID)]; !ok {
		return fmt.Errorf("user not found")
	}
	password := r.storage[fmt.Sprintf("%d", userID)]["password"]
	r.storage[password]["email_verified"] = "true"
	return nil
}

func (r *MemoryRepository) UpdateDisplayName(userID uint, req dto.UpdateDisplayNameDTO) error {
	if _, ok := r.storage[fmt.Sprintf("%d", userID)]; !ok {
		return fmt.Errorf("user not found")
	}
	password := r.storage[fmt.Sprintf("%d", userID)]["password"]
	r.storage[password]["display_name"] = req.DisplayName
	return nil
}

// Oauth methods
func (r *MemoryRepository) LoginUserPrecheck(username string) (string, error) {
	if _, ok := r.storage[username]; !ok {
		fmt.Println("user not found while using LoginUserPrecheck")
		return "", fmt.Errorf("user not found")
	}
	return r.storage[username]["salt"], nil
}

func (r *MemoryRepository) LoginPreCheckClient(req dto.LoginPrecheckDTO) (string, string, error) {

	return "r.storage[username]['salt']", "second string...", nil
}

func (r *MemoryRepository) GetUser(username string) (*serverModels.User, error) {
	if _, ok := r.storage[username]; !ok {
		fmt.Println("user not found while using GetUser")
		return &serverModels.User{}, fmt.Errorf("user not found")
	}
	password := r.storage[username]["password"]
	userID := r.storage[password]["user_id"]
	userIdInt, err := strconv.Atoi(userID)
	if err != nil {
		return &serverModels.User{}, err
	}
	user := serverModels.User{
		ID:        uint(userIdInt),
		Email:     r.storage[password]["email"],
		Username:  r.storage[password]["username"],
		FirstName: r.storage[password]["first_name"],
		LastName:  r.storage[password]["last_name"],
		Password:  r.storage[password]["password"],
		Salt:      r.storage[username]["salt"],
	}
	return &user, nil
}

func (r *MemoryRepository) GetUserByID(id int64) (*serverModels.User, error) {
	if _, ok := r.storage[fmt.Sprintf("%d", id)]; !ok {
		fmt.Println("user not found while using GetUserByID")
		return &serverModels.User{}, fmt.Errorf("user not found")
	}
	password := r.storage[fmt.Sprintf("%d", id)]["password"]
	user := serverModels.User{
		ID:        uint(id),
		Email:     r.storage[password]["email"],
		Username:  r.storage[password]["username"],
		FirstName: r.storage[password]["first_name"],
		LastName:  r.storage[password]["last_name"],
		Password:  r.storage[password]["password"],
		Salt:      r.storage[password]["salt"],
	}
	return &user, nil
}

func (r *MemoryRepository) GetUserMetadata(userID int64, key string) (*serverModels.UserMetadata, error) {
	if _, ok := r.storage[fmt.Sprintf("%d", userID)]; !ok {
		fmt.Println("user not found while using GetUserMetadata")
		return &serverModels.UserMetadata{}, fmt.Errorf("user not found")
	}
	password := r.storage[fmt.Sprintf("%d", userID)]["password"]
	userMetadata := serverModels.UserMetadata{
		Key:   key,
		Value: r.storage[password][key],
	}
	return &userMetadata, nil
}

func (r *MemoryRepository) SetClient(client *serverModels.Client) error {
	r.storage[client.ID] = map[string]string{
		"id":           client.ID,
		"secret":       client.Secret,
		"name":         client.Name,
		"redirect_uri": client.RedirectURI,
	}
	return nil
}

func (r *MemoryRepository) GetClient(id string) (*serverModels.Client, error) {
	if strings.Contains(id, ":") {
		id = id[strings.LastIndex(id, ":")+1:]
		// fmt.Println("ID check:", id)
	}
	if _, ok := r.storage[id]; !ok {
		fmt.Println("client not found while using GetClient")
		return &serverModels.Client{}, fmt.Errorf("client not found")
	}
	client := serverModels.Client{
		ID:          r.storage[id]["id"],
		Secret:      r.storage[id]["secret"],
		Name:        r.storage[id]["name"],
		RedirectURI: r.storage[id]["redirect_uri"],
	}
	return &client, nil
}

func (r *MemoryRepository) SetTTL(key string, value []byte, ttl time.Duration) error {
	r.byteStorage[key] = value
	go func() {
		time.Sleep(ttl)
		delete(r.storage, key)
	}()
	return nil
}

func (r *MemoryRepository) GetTTL(key string) ([]byte, error) {
	if _, ok := r.byteStorage[key]; !ok {
		fmt.Println("key not found while using GetTTL")
		return nil, fmt.Errorf("key not found")
	}
	return r.byteStorage[key], nil
}
