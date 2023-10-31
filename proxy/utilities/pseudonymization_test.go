package utilities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePlaceholderUserData(t *testing.T) {
	username, email, fname, lname := GeneratePlaceholderUserData()
	assert.NotEmpty(t, username)
	assert.NotEmpty(t, email)
	assert.NotEmpty(t, fname)
	assert.NotEmpty(t, lname)
}
