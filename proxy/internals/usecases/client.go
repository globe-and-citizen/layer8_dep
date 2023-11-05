package usecases

import (
	"encoding/json"
	"fmt"
	"globe-and-citizen/layer8/proxy/constants"
	"globe-and-citizen/layer8/proxy/entities"

	"github.com/google/uuid"
)

func (u *UseCase) GetClient(id string) (*entities.Client, error) {
	var client entities.Client
	res := u.Repo.Get(fmt.Sprintf("client:%s", id))
	if res == nil {
		return nil, constants.ErrNotFound
	}
	err := json.Unmarshal(res, &client)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (u *UseCase) AddClient(name string) (*entities.Client, error) {
	client := &entities.Client{
		ID:     uuid.New().String()[:8],
		Secret: uuid.New().String()[:8],
		Name:   name,
	}
	b, err := json.Marshal(client)
	if err != nil {
		return nil, err
	}
	err = u.Repo.Set(fmt.Sprintf("client:%s", client.ID), b)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// this is only be used for testing purposes
func (u *UseCase) AddTestClient() (*entities.Client, error) {
	client := &entities.Client{
		ID:          "notanid",
		Secret:      "absolutelynotasecret!",
		Name:        "Ex-C",
		RedirectURI: "http://localhost:5173/oauth2/callback",
	}
	b, err := json.Marshal(client)
	if err != nil {
		return nil, err
	}
	err = u.Repo.Set(fmt.Sprintf("client:%s", client.ID), b)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (u *UseCase) DeleteClient(id string) error {
	return u.Repo.Delete(fmt.Sprintf("client:%s", id))
}
