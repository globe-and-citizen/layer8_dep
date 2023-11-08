package usecases

import (
	"encoding/json"
	"fmt"
	"globe-and-citizen/layer8/proxy/models"
)

func (u *UseCase) GetClient(id string) (*models.Client, error) {
	client := u.Repo.GetClient(fmt.Sprintf("client:%s", id))

	err := json.Unmarshal(res, &client)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// this is only be used for testing purposes
func (u *UseCase) AddTestClient() (*models.Client, error) {
	client := &models.Client{
		ID:          "notanid",
		Secret:      "absolutelynotasecret!",
		Name:        "Ex-C",
		RedirectURI: "http://localhost:5173/oauth2/callback",
	}
	err := u.Repo.SetClient(client)
	if err != nil {
		return nil, err
	}
	return client, nil
}
