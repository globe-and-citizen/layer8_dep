package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/interfaces"
	"globe-and-citizen/layer8/server/resource_server/utils"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	utils.GetPwd()

	// var relativePathIndex = "assets-v1/templates_copy/homeView.html"
	var relativePathIndex = "assets-v1/templates_copy/public/welcome.html"
	indexPath := filepath.Join(utils.WorkingDirectory, relativePathIndex)
	fmt.Println("indexPath: ", indexPath)
	http.ServeFile(w, r, indexPath)
}
func LoginUserPage(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	utils.GetPwd()

	var relativePathIndex = "assets-v1/templates_copy/src/pages/user_portal/login.html"
	indexPath := filepath.Join(utils.WorkingDirectory, relativePathIndex)
	http.ServeFile(w, r, indexPath)
}
func RegisterUserPage(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	utils.GetPwd()

	var relativePathIndex = "assets-v1/templates_copy/src/pages/user_portal/register.html"
	indexPath := filepath.Join(utils.WorkingDirectory, relativePathIndex)
	http.ServeFile(w, r, indexPath)
}

func ClientProfilePage(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	utils.GetPwd()

	var relativePathUser = "assets-v1/templates_copy/src/pages/client_portal/profile.html"
	userPath := filepath.Join(utils.WorkingDirectory, relativePathUser)
	fmt.Println("userPath: ", userPath)
	http.ServeFile(w, r, userPath)
}
func UserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	utils.GetPwd()

	var relativePathUser = "assets-v1/templates_copy/src/pages/user_portal/profile.html"
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

	var relativePathUser = "assets-v1/templates_copy/src/pages/client_portal/register.html"
	userPath := filepath.Join(utils.WorkingDirectory, relativePathUser)
	fmt.Println("userPath: ", userPath)
	http.ServeFile(w, r, userPath)
}
func LoginClientPage(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	utils.GetPwd()

	var relativePathUser = "assets-v1/templates_copy/src/pages/client_portal/login.html"
	userPath := filepath.Join(utils.WorkingDirectory, relativePathUser)
	fmt.Println("userPath: ", userPath)
	http.ServeFile(w, r, userPath)
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	newService := r.Context().Value("service").(interfaces.IService)
	var req dto.RegisterUserDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}

	err := newService.RegisterUser(req)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}

	res := utils.BuildResponse(true, "OK!", "User registered successfully")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}
}

func RegisterClientHandler(w http.ResponseWriter, r *http.Request) {
	newService := r.Context().Value("service").(interfaces.IService)
	var req dto.RegisterClientDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to register client", err)
		return
	}

	err := newService.RegisterClient(req)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}

	res := utils.BuildResponse(true, "OK!", "Client registered successfully")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to register client", err)
		return
	}
}

// LoginPrecheckHandler handles login precheck requests and get the salt of the user from the database using the username from the request URL
func LoginPrecheckHandler(w http.ResponseWriter, r *http.Request) {
	newService := r.Context().Value("service").(interfaces.IService)
	var req dto.LoginPrecheckDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}

	loginPrecheckResp, err := newService.LoginPreCheckUser(req)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}

	if err := json.NewEncoder(w).Encode(loginPrecheckResp); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}
}
func ClientLoginPrecheckHandler(w http.ResponseWriter, r *http.Request) {
	newService := r.Context().Value("service").(interfaces.IService)
	var req dto.LoginPrecheckDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}

	loginPrecheckResp, err := newService.LoginPreCheckClient(req)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}

	if err := json.NewEncoder(w).Encode(loginPrecheckResp); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	newService := r.Context().Value("service").(interfaces.IService)
	var req dto.LoginUserDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}

	tokenResp, err := newService.LoginUser(req)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}

	if err := json.NewEncoder(w).Encode(tokenResp); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}
}

func LoginClientHandler(w http.ResponseWriter, r *http.Request) {
	newService := r.Context().Value("service").(interfaces.IService)
	var req dto.LoginClientDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}

	tokenResp, err := newService.LoginClient(req)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}

	if err := json.NewEncoder(w).Encode(tokenResp); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	newService := r.Context().Value("service").(interfaces.IService)
	tokenString := r.Header.Get("Authorization")
	tokenString = tokenString[7:] // Remove the "Bearer " prefix
	userID, err := utils.ValidateToken(tokenString)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get user profile", err)
		return
	}

	profileResp, err := newService.ProfileUser(userID)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get user profile", err)
		return
	}

	if err := json.NewEncoder(w).Encode(profileResp); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get user profile", err)
		return
	}
}

func ClientProfileHandler(w http.ResponseWriter, r *http.Request) {
	newService := r.Context().Value("service").(interfaces.IService)
	tokenString := r.Header.Get("Authorization")
	tokenString = tokenString[7:]
	userID, err := utils.ValidateClientToken(tokenString)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get user profile", err)
		return
	}

	profileResp, err := newService.ProfileClient(userID)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get user profile", err)
		return
	}

	if err := json.NewEncoder(w).Encode(profileResp); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get user profile", err)
		return
	}
}

func GetClientData(w http.ResponseWriter, r *http.Request) {
	newService := r.Context().Value("service").(interfaces.IService)
	clientName := r.Header.Get("Name")

	clientModel, err := newService.GetClientData(clientName)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}

	if err := json.NewEncoder(w).Encode(clientModel); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}
}

func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	newService := r.Context().Value("service").(interfaces.IService)
	tokenString := r.Header.Get("Authorization")
	tokenString = tokenString[7:] // Remove the "Bearer " prefix
	userID, err := utils.ValidateToken(tokenString)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to verify email", err)
		return
	}

	err = newService.VerifyEmail(userID)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to verify email", err)
		return
	}

	resp := utils.BuildResponse(true, "OK!", "Email verified successfully")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to verify email", err)
		return
	}
}

func UpdateDisplayNameHandler(w http.ResponseWriter, r *http.Request) {
	newService := r.Context().Value("service").(interfaces.IService)
	tokenString := r.Header.Get("Authorization")
	tokenString = tokenString[7:] // Remove the "Bearer " prefix
	userID, err := utils.ValidateToken(tokenString)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to update display name", err)
		return
	}

	var req dto.UpdateDisplayNameDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to update display name", err)
		return
	}

	err = newService.UpdateDisplayName(userID, req)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to update display name", err)
		return
	}

	resp := utils.BuildResponse(true, "OK!", "Display name updated successfully")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to update display name", err)
		return
	}
}
