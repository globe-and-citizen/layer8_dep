package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"globe-and-citizen/layer8/server/config"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/interfaces"
	"globe-and-citizen/layer8/server/resource_server/models"
	"globe-and-citizen/layer8/server/resource_server/utils"

	"github.com/go-playground/validator/v10"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	utils.GetPwd()

	var relativePathIndex = "assets-v1/templates/homeView.html"
	indexPath := filepath.Join(utils.WorkingDirectory, relativePathIndex)
	fmt.Println("indexPath: ", indexPath)
	http.ServeFile(w, r, indexPath)

}

func UserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	utils.GetPwd()

	var relativePathUser = "assets-v1/templates/userView.html"
	userPath := filepath.Join(utils.WorkingDirectory, relativePathUser)
	fmt.Println("userPath: ", userPath)
	http.ServeFile(w, r, userPath)
}
func ClientHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	utils.GetPwd()

	var relativePathUser = "assets-v1/templates/registerClient.html"
	userPath := filepath.Join(utils.WorkingDirectory, relativePathUser)
	fmt.Println("userPath: ", userPath)
	http.ServeFile(w, r, userPath)
}

// RegisterUserHandler handles user registration requests
func RegisterUserHandler(w http.ResponseWriter, r *http.Request, service interfaces.IService) {

	var req dto.RegisterUserDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := utils.BuildErrorResponse("Failed to register user", err.Error(), utils.EmptyObj{})
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}

	err := service.RegisterUser(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := utils.BuildErrorResponse("Failed to register user", err.Error(), utils.EmptyObj{})
		json.NewEncoder(w).Encode(res)
		return
	}

	res := utils.BuildResponse(true, "OK!", "User registered successfully")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := utils.BuildErrorResponse("Failed to register user", err.Error(), utils.EmptyObj{})
		json.NewEncoder(w).Encode(res)
		return
	}
}

// Register Client
func RegisterClientHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterClientDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, http.StatusBadRequest, "Failed to register client", err)
		return
	}

	if err := validator.New().Struct(req); err != nil {
		handleError(w, http.StatusBadRequest, "Failed to register client", err)
		return
	}

	clientUUID := utils.GenerateUUID()
	clientSecret := utils.GenerateSecret(utils.SecretSize)

	// db := config.SetupDatabaseConnection()
	// defer config.CloseDatabaseConnection(db)

	client := models.Client{
		ID:          clientUUID,
		Secret:      clientSecret,
		Name:        req.Name,
		RedirectURI: req.RedirectURI,
	}

	if err := config.DB.Create(&client).Error; err != nil {
		handleError(w, http.StatusBadRequest, "Failed to register client", err)
		return
	}

	res := utils.BuildResponse(true, "OK!", "Client registered successfully")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		handleError(w, http.StatusBadRequest, "Failed to register client", err)
		return
	}
}

// LoginPrecheckHandler handles login precheck requests and get the salt of the user from the database using the username from the request URL
func LoginPrecheckHandler(w http.ResponseWriter, r *http.Request) {

	var req dto.LoginPrecheckDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := utils.BuildErrorResponse("Failed to perform login precheck", err.Error(), utils.EmptyObj{})
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}

	// validate request
	if err := validator.New().Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := utils.BuildErrorResponse("Failed to perform login precheck", err.Error(), utils.EmptyObj{})
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	// Using the username, find the user in the database
	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := utils.BuildErrorResponse("Failed to perform login precheck", err.Error(), utils.EmptyObj{})
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	resp := models.LoginPrecheckResponseOutput{
		Username: user.Username,
		Salt:     user.Salt,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := utils.BuildErrorResponse("Failed to perform login precheck", err.Error(), utils.EmptyObj{})
		json.NewEncoder(w).Encode(res)
		return
	}
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {

	var req dto.LoginUserDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := utils.BuildErrorResponse("Failed to login user", err.Error(), utils.EmptyObj{})
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}

	// validate request
	if err := validator.New().Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := utils.BuildErrorResponse("Failed to login user", err.Error(), utils.EmptyObj{})
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	// Using the username, find the user in the database
	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := utils.BuildErrorResponse("Failed to login user", err.Error(), utils.EmptyObj{})
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}

	HashedAndSaltedPass := utils.SaltAndHashPassword(req.Password, user.Salt)

	// Compare the password with the password in the database
	if user.Password != HashedAndSaltedPass {
		w.WriteHeader(http.StatusUnauthorized)
		res := utils.BuildErrorResponse("Error", "Invalid credentials", utils.EmptyObj{})
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}

	JWT_SECRET_STR := os.Getenv("JWT_SECRET")
	JWT_SECRET_BYTE := []byte(JWT_SECRET_STR)

	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &models.Claims{
		UserName: user.Username,
		UserID:   user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    "GlobeAndCitizen",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWT_SECRET_BYTE)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := utils.BuildErrorResponse("Failed to login user", err.Error(), utils.EmptyObj{})
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	resp := models.LoginUserResponseOutput{
		Token: tokenString,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := utils.BuildErrorResponse("Failed to login user", err.Error(), utils.EmptyObj{})
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {

	// Get the token from the request header
	tokenString := r.Header.Get("Authorization")
	tokenString = tokenString[7:] // Remove the "Bearer " prefix
	// Get user ID from the token
	claims := &models.Claims{}
	JWT_SECRET_STR := os.Getenv("JWT_SECRET")
	JWT_SECRET_BYTE := []byte(JWT_SECRET_STR)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JWT_SECRET_BYTE, nil
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := utils.BuildErrorResponse("Failed to get user profile", err.Error(), utils.EmptyObj{})
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		res := utils.BuildErrorResponse("Failed to get user profile", err.Error(), utils.EmptyObj{})
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}

	var user models.User
	if err := config.DB.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := utils.BuildErrorResponse("Failed to get user profile", err.Error(), utils.EmptyObj{})
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}

	// Get user metadata from the database
	var userMetadata []models.UserMetadata
	if err := config.DB.Where("user_id = ?", claims.UserID).Find(&userMetadata).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := utils.BuildErrorResponse("Failed to get user profile", err.Error(), utils.EmptyObj{})
		json.NewEncoder(w).Encode(res)
		return
	}

	resp := models.ProfileResponseOutput{
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	for _, metadata := range userMetadata {
		switch metadata.Key {
		case "display_name":
			resp.DisplayName = metadata.Value
		case "country":
			resp.Country = metadata.Value
		case "email_verified":
			resp.EmailVerified = metadata.Value == "true"
		}
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := utils.BuildErrorResponse("Failed to get user profile", err.Error(), utils.EmptyObj{})
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
}
func GetClientData(w http.ResponseWriter, r *http.Request) {
	clientName := r.Header.Get("Name")

	// db := config.SetupDatabaseConnection()
	// defer config.CloseDatabaseConnection(db)

	var client models.Client
	if err := config.DB.Where("name = ?", clientName).First(&client).Error; err != nil {
		handleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}

	resp := models.ClientResponseOutput{
		ID:          client.ID,
		Secret:      client.Secret,
		Name:        client.Name,
		RedirectURI: client.RedirectURI,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		handleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}
}

func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {

	// Get the token from the request header
	tokenString := r.Header.Get("Authorization")
	tokenString = tokenString[7:] // Remove the "Bearer " prefix
	// Get user ID from the token
	claims := &models.Claims{}
	JWT_SECRET_STR := os.Getenv("JWT_SECRET")
	JWT_SECRET_BYTE := []byte(JWT_SECRET_STR)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JWT_SECRET_BYTE, nil
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := utils.BuildErrorResponse("Failed to verify email", err.Error(), utils.EmptyObj{})
		json.NewEncoder(w).Encode(res)
		return
	}
	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		res := utils.BuildErrorResponse("Failed to verify email", err.Error(), utils.EmptyObj{})
		json.NewEncoder(w).Encode(res)
		return
	}

	err = config.DB.Model(&models.UserMetadata{}).Where("user_id = ? AND key = ?", claims.UserID, "email_verified").Update("value", "true").Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := utils.BuildErrorResponse("Failed to verify email", err.Error(), utils.EmptyObj{})
		json.NewEncoder(w).Encode(res)
		return
	}

	resp := utils.BuildResponse(true, "OK!", "Email verified successfully")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := utils.BuildErrorResponse("Failed to verify email", err.Error(), utils.EmptyObj{})
		json.NewEncoder(w).Encode(res)
		return
	}
}

func UpdateDisplayNameHandler(w http.ResponseWriter, r *http.Request) {

	var req dto.UpdateDisplayNameDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := utils.BuildErrorResponse("Failed to update display name", err.Error(), utils.EmptyObj{})
		json.NewEncoder(w).Encode(res)
		return
	}

	// validate request
	if err := validator.New().Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := utils.BuildErrorResponse("Failed to update display name", err.Error(), utils.EmptyObj{})
		json.NewEncoder(w).Encode(res)
		return
	}

	// Get the token from the request header
	tokenString := r.Header.Get("Authorization")
	tokenString = tokenString[7:] // Remove the "Bearer " prefix
	// Get user ID from the token
	claims := &models.Claims{}
	JWT_SECRET_STR := os.Getenv("JWT_SECRET")
	JWT_SECRET_BYTE := []byte(JWT_SECRET_STR)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JWT_SECRET_BYTE, nil
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := utils.BuildErrorResponse("Failed to update display name", err.Error(), utils.EmptyObj{})
		json.NewEncoder(w).Encode(res)
		return
	}
	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		res := utils.BuildErrorResponse("Failed to update display name", err.Error(), utils.EmptyObj{})
		json.NewEncoder(w).Encode(res)
		return
	}

	err = config.DB.Model(&models.UserMetadata{}).Where("user_id = ? AND key = ?", claims.UserID, "display_name").Update("value", req.DisplayName).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := utils.BuildErrorResponse("Failed to update display name", err.Error(), utils.EmptyObj{})
		json.NewEncoder(w).Encode(res)
		return
	}

	resp := utils.BuildResponse(true, "OK!", "Display name updated successfully")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := utils.BuildErrorResponse("Failed to update display name", err.Error(), utils.EmptyObj{})
		json.NewEncoder(w).Encode(res)
		return
	}
}

func handleError(w http.ResponseWriter, status int, message string, err error) {
	w.WriteHeader(status)
	res := utils.BuildErrorResponse(message, err.Error(), utils.EmptyObj{})
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Error sending response: %v", err)
	}
}
