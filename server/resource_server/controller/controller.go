package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"globe-and-citizen/layer8/proxy/resource_server/config"
	"globe-and-citizen/layer8/proxy/resource_server/dto"
	"globe-and-citizen/layer8/proxy/resource_server/models"
	"globe-and-citizen/layer8/proxy/resource_server/utils"

	"github.com/go-playground/validator/v10"
)

// RegisterUserHandler handles user registration requests
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {

	var req dto.RegisterUserDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := utils.BuildErrorResponse("Failed to register user", err.Error(), utils.EmptyObj{})
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}

	// validate request
	if err := validator.New().Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := utils.BuildErrorResponse("Failed to register user", err.Error(), utils.EmptyObj{})
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}

	rmSalt := utils.GenerateRandomSalt(utils.SaltSize)
	HashedAndSaltedPass := utils.SaltAndHashPassword(req.Password, rmSalt)

	// Make connection to database
	db := config.SetupDatabaseConnection()
	// Close connection database
	defer config.CloseDatabaseConnection(db)
	// Save user to database
	user := models.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  HashedAndSaltedPass,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Salt:      rmSalt,
	}
	if err := db.Create(&user).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := utils.BuildErrorResponse("Failed to register user", err.Error(), utils.EmptyObj{})
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
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
	if err := db.Create(&userMetadata).Error; err != nil {
		db.Delete(&user)
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
	// Make connection to database
	db := config.SetupDatabaseConnection()
	// Close connection database
	defer config.CloseDatabaseConnection(db)
	// Using the username, find the user in the database
	var user models.User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
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
	// Make connection to database
	db := config.SetupDatabaseConnection()
	// Close connection database
	defer config.CloseDatabaseConnection(db)
	// Using the username, find the user in the database
	var user models.User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
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

	db := config.SetupDatabaseConnection()
	defer config.CloseDatabaseConnection(db)
	var user models.User
	if err := db.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := utils.BuildErrorResponse("Failed to get user profile", err.Error(), utils.EmptyObj{})
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}

	// Get user metadata from the database
	var userMetadata []models.UserMetadata
	if err := db.Where("user_id = ?", claims.UserID).Find(&userMetadata).Error; err != nil {
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

	db := config.SetupDatabaseConnection()
	defer config.CloseDatabaseConnection(db)

	err = db.Model(&models.UserMetadata{}).Where("user_id = ? AND key = ?", claims.UserID, "email_verified").Update("value", "true").Error
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

	db := config.SetupDatabaseConnection()
	defer config.CloseDatabaseConnection(db)

	err = db.Model(&models.UserMetadata{}).Where("user_id = ? AND key = ?", claims.UserID, "display_name").Update("value", req.DisplayName).Error
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
