package utilities

import (
	"strings"
	"time"

	"github.com/goombaio/namegenerator"
)

// GeneratePlaceholderUserData generates placeholder user data.
// This is used to pseudonymize the user data.
// Returns the username, email, first name and last name.
func GeneratePlaceholderUserData() (string, string, string, string) {
	seed := time.Now().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	username := strings.ToLower(nameGenerator.Generate())
	email := strings.ReplaceAll(username, "-", "") + "@placeholder.com"
	fullname := strings.Split(username, "-")
	fname := strings.ToUpper(fullname[0][:1]) + fullname[0][1:]
	lname := strings.ToUpper(fullname[1][:1]) + fullname[1][1:]
	return username, email, fname, lname
}
