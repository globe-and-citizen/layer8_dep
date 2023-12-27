package dto

type RegisterUserDTO struct {
	Email       string `json:"email" validate:"required,email"`
	Username    string `json:"username" validate:"required,min=3,max=50"`
	Password    string `json:"password" validate:"required"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name"`
	DisplayName string `json:"display_name"`
	Country     string `json:"country"`
}

type LoginUserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

type LoginPrecheckDTO struct {
	Username string `json:"username"`
}

type UpdateDisplayNameDTO struct {
	DisplayName string `json:"display_name"`
}
