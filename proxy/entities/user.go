package entities

import (
	"globe-and-citizen/layer8/l8_oauth/utilities"
)

type AbstractUser struct {
	Username                     string `json:"username" validate:"nonzero"`
	Fname                        string `json:"first_name" validate:"nonzero"`
	Lname                        string `json:"last_name" validate:"nonzero"`
	Email                        string `json:"email" validate:"nonzero"`
	PhoneNumber                  string `json:"phone_number" validate:"nonzero"`
	Address                      string `json:"address" validate:"nonzero"`
	NationalIdentificationNumber string `json:"national_identification_number" validate:"nonzero"`
	ShareEmailVer                bool   `json:"share_email_ver"`
	SharePhoneNumberVer          bool   `json:"share_phone_number_ver"`
	ShareAddressVer              bool   `json:"share_address_ver"`
	ShareIdVer                   bool   `json:"share_id_ver"`
}

type User struct {
	ID               int64        `json:"id"`
	Password         string       `json:"password" validate:"nonzero"`
	PsedonymizedData AbstractUser `json:"pseudonymized_data"`
	AbstractUser
}

type ReqUserData struct {
	Username            string `json:"username"`
	ShareEmailVer       bool   `json:"share_email_ver"`
	SharePhoneNumberVer bool   `json:"share_phone_number_ver"`
	ShareAddressVer     bool   `json:"share_address_ver"`
	ShareIdVer          bool   `json:"share_id_ver"`
}

// Validate validates the user struct by checking if the required fields are
// present and if the email is valid.
func (u *User) Validate() error {
	// Check if the required fields are present
	err := utilities.ValidateRequiredFields(u)
	if err != nil {
		return err
	}

	// Check if the email is valid
	err = utilities.ValidateEmail(u.Email)
	if err != nil {
		return err
	}

	return nil
}
