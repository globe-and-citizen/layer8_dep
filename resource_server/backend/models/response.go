package models

type LoginPrecheckResponseOutput struct {
	Username string `json:"username"`
	Salt     string `json:"salt"`
}

type LoginUserResponseOutput struct {
	Token string `json:"token"`
}

type ProfileResponseOutput struct {
	Email               string `json:"email"`
	Username            string `json:"username"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	PhoneNumber         string `json:"phone_number"`
	Address             string `json:"address"`
	EmailVerified       bool   `json:"email_verified"`
	PhoneNumberVerified bool   `json:"phone_number_verified"`
	LocationVerified    bool   `json:"location_verified"`
	NationalIdVerified  bool   `json:"national_id_verified"`
}

type ExposeUserResponseOutput struct {
	Username            string `json:"username"`
	EmailVerified       bool   `json:"email_verified"`
	PhoneNumberVerified bool   `json:"phone_number_verified"`
	LocationVerified    bool   `json:"location_verified"`
	NationalIdVerified  bool   `json:"national_id_verified"`
}
