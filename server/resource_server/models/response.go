package models

type LoginPrecheckResponseOutput struct {
	Username string `json:"username"`
	Salt     string `json:"salt"`
}

type LoginUserResponseOutput struct {
	Token string `json:"token"`
}

type ProfileResponseOutput struct {
	Email         string `json:"email"`
	Username      string `json:"username"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	DisplayName   string `json:"display_name"`
	Country       string `json:"country"`
	EmailVerified bool   `json:"email_verified"`
}
