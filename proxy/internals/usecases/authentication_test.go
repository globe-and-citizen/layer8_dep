package usecases

import (
	"globe-and-citizen/layer8/proxy/entities"
	"globe-and-citizen/layer8/proxy/internals/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterLogingGetUser(t *testing.T) {
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

	// register user
	res, err := usecase.RegisterUser(&entities.User{
		AbstractUser: entities.AbstractUser{
			Username: "test",
			Email:    "test@mail.com",
			Fname:    "test fname",
			Lname:    "test lname",
		},
		Password: "test",
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)

	token := res["token"].(string)
	user := res["user"].(*entities.User)
	assert.NotEmpty(t, token)
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

	// login user
	res, err = usecase.LoginUser(user.Username, "wrong password")
	assert.Error(t, err)
	assert.Nil(t, res)

	res, err = usecase.LoginUser(user.Username, "test")
	assert.NoError(t, err)
	assert.NotNil(t, res)

	// get user by token
	token = res["token"].(string)
	user = res["user"].(*entities.User)
	assert.NotEmpty(t, token)
	assert.NotNil(t, user)
	user2, err := usecase.GetUserByToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, user2)
	assert.Equal(t, user.ID, user2.ID)
	assert.Equal(t, user.Username, user2.Username)
	assert.Equal(t, user.Email, user2.Email)
	assert.Equal(t, user.Fname, user2.Fname)
	assert.Equal(t, user.Lname, user2.Lname)

	// delete user
	err = usecase.DeleteUser(user.ID)
	assert.NoError(t, err)
}
