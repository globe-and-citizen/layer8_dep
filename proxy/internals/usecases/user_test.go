package usecases

import (
	"testing"

	"globe-and-citizen/layer8/l8_oauth/entities"
	"globe-and-citizen/layer8/l8_oauth/internals/repository"

	"github.com/stretchr/testify/assert"
)

func TestAddGetDeleteUSer(t *testing.T) {
	usecase := &UseCase{Repo: repository.MustCreateRepository("memory")}

	// check checks that all fields are set
	check := func(u *entities.User) {
		assert.NotEmpty(t, u.ID)
		assert.Equal(t, "test", u.Username)
		assert.Equal(t, "test@mail.com", u.Email)
		assert.Equal(t, "test fname", u.Fname)
		assert.Equal(t, "test lname", u.Lname)
		assert.NotEqual(t, "test", u.Password)
		assert.NotEmpty(t, u.PsedonymizedData.Username)
		assert.NotEmpty(t, u.PsedonymizedData.Email)
		assert.NotEmpty(t, u.PsedonymizedData.Fname)
		assert.NotEmpty(t, u.PsedonymizedData.Lname)
	}

	user, err := usecase.AddUser(&entities.User{
		AbstractUser: entities.AbstractUser{
			Username: "test",
			Email:    "test@mail.com",
			Fname:    "test fname",
			Lname:    "test lname",
		},
		Password: "test",
	})
	assert.NoError(t, err)
	assert.NotNil(t, user)
	check(user)

	// get original user data
	user, err = usecase.GetUser(user.ID, false)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	check(user)

	// get pseudonymized user data
	puser, err := usecase.GetUser(user.ID, true)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	// check that the pseudonymized data is not the same as the original data
	assert.NotEqual(t, user.Username, puser.Username)
	assert.NotEqual(t, user.Email, puser.Email)
	assert.NotEqual(t, user.Fname, puser.Fname)
	assert.NotEqual(t, user.Lname, puser.Lname)

	// delete user
	err = usecase.DeleteUser(user.ID)
	assert.NoError(t, err)
}
