package usecases

import (
	"testing"

	"globe-and-citizen/layer8/l8_oauth/internals/repository"

	"github.com/stretchr/testify/assert"
)

func TestAddGetDeleteClient(t *testing.T) {
	usecase := &UseCase{Repo: repository.MustCreateRepository("memory")}

	client, err := usecase.AddClient("test")
	assert.NoError(t, err)
	assert.NotNil(t, client)

	client, err = usecase.GetClient(client.ID)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	err = usecase.DeleteClient(client.ID)
	assert.NoError(t, err)

	client, err = usecase.GetClient(client.ID)
	assert.Error(t, err)
	assert.Nil(t, client)
}
