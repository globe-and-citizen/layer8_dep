package usecases

import (
	"layer8-proxy/constants"
	"layer8-proxy/entities"
	"layer8-proxy/internals/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestGenerateAuthorizationURL_ExchangeCodeForToken_AccessResourcesWithToken(t *testing.T) {
	var (
		usecase = &UseCase{Repo: repository.MustCreateRepository("memory")}
		client  = &entities.Client{
			Name: "Test Client",
		}
		config = &oauth2.Config{
			RedirectURL: "http://localhost:8080",
		}
	)

	// client and user not registered
	url, err := usecase.GenerateAuthorizationURL(config, 1)
	assert.Error(t, err)
	assert.Nil(t, url)

	// register client
	client, err = usecase.AddClient(client.Name)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	config.ClientID = client.ID
	config.ClientSecret = client.Secret

	// register user
	user, err := usecase.AddUser(&entities.User{
		AbstractUser: entities.AbstractUser{
			Username: "test",
			Email:    "test@mail.com",
			Fname:    "test fname",
			Lname:    "test lname",
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, user)

	// both user and client are now registered
	// generate authorization url (with no scopes)
	url, err = usecase.GenerateAuthorizationURL(config, user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, url)

	// exchange code for token (with no scopes)
	token, err := usecase.ExchangeCodeForToken(config, url.Code)
	assert.Nil(t, err)
	assert.NotNil(t, token)

	// get data with token (with no scopes)
	resource, err := usecase.AccessResourcesWithToken(token.AccessToken)
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{}, resource)

	// generate authorization url (with scopes)
	config.Scopes = []string{constants.READ_USER_SCOPE}
	url, err = usecase.GenerateAuthorizationURL(config, user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, url)

	// exchange code for token (with scopes)
	token, err = usecase.ExchangeCodeForToken(config, url.Code)
	assert.Nil(t, err)
	assert.NotNil(t, token)

	// get data with token (with scopes)
	resource, err = usecase.AccessResourcesWithToken(token.AccessToken)
	assert.NoError(t, err)
	assert.NotNil(t, resource)
	user, ok := resource["profile"].(*entities.User)
	assert.True(t, ok)
	assert.NotEqual(t, "test", user.Username)
	assert.NotEqual(t, "test@mail.com", user.Email)
	assert.NotEqual(t, "test fname", user.Fname)
	assert.NotEqual(t, "test lname", user.Lname)

	// cleanup
	err = usecase.DeleteClient(client.ID)
	assert.NoError(t, err)
	err = usecase.DeleteUser(user.ID)
	assert.NoError(t, err)
}
