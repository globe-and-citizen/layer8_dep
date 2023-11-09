package dto

type RegisterUserDTO struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"`
}

type LoginUserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

type LoginPrecheckDTO struct {
	Username string `json:"username"`
}
