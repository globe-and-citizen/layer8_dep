package utilities

import (
	"log"
	"regexp"

	"globe-and-citizen/layer8/l8_oauth/constants"

	"gopkg.in/validator.v2"
)

// ValidateRequiredFields validates if the required fields are present in the struct.
func ValidateRequiredFields(s interface{}) error {
	err := validator.Validate(s)
	if err != nil {
		log.Println(err.(validator.ErrorMap))
		return constants.ErrMissingFields
	}
	return nil
}

// ValidateEmail validates if the email is valid.
func ValidateEmail(email string) error {
	regex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !regex.MatchString(email) {
		return constants.ErrInvalidEmail
	}
	return nil
}
