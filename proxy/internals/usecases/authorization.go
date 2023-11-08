package usecases

import (
	"encoding/json"
	"fmt"
	"globe-and-citizen/layer8/proxy/constants"
	"globe-and-citizen/layer8/proxy/entities"
	"globe-and-citizen/layer8/utils"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

// GenerateAuthorizationURL generates an authorization URL for the user to visit
// and authorize the application to access their account.
func (u *UseCase) GenerateAuthorizationURL(config *oauth2.Config, userID int64) (*entities.AuthURL, error) {
	// first, check that both client and user exist
	client, err := u.GetClient(config.ClientID)
	if err != nil {
		return nil, fmt.Errorf("could not get client: %v", err)
	}
	user, err := u.GetUser(userID, true)
	if err != nil {
		return nil, fmt.Errorf("could not get user: %v", err)
	}

	state, stateErr := utils.GenerateRandomString(24)
	if stateErr != nil {
		return nil, fmt.Errorf("could not generate random state: %v", stateErr)
	}

	// generate the auth code
	scopes := ""
	for _, scope := range config.Scopes {
		scopes += scope + ","
	}
	code, err := utils.GenerateAuthCode(client.Secret, &utils.AuthCodeClaims{
		ClientID:    config.ClientID,
		UserID:      user.ID,
		RedirectURI: config.RedirectURL,
		Scopes:      scopes,
		ExpiresAt:   time.Now().Add(time.Minute * 5).Unix(),
	})
	if err != nil {
		return nil, fmt.Errorf("could not generate auth code: %v", err)
	}

	return &entities.AuthURL{
		URL: config.AuthCodeURL(
			state,
			oauth2.SetAuthURLParam("code", code),
		),
		Code:  code,
		State: state,
	}, nil
}

// ExchangeCodeForToken generates an access token from an authorization code.
func (u *UseCase) ExchangeCodeForToken(config *oauth2.Config, code string) (*oauth2.Token, error) {
	// ensure that the secret is specified
	if config.ClientSecret == "" {
		return nil, fmt.Errorf("client secret is not specified")
	}
	// verify the code
	claims, err := utils.DecodeAuthCode(config.ClientSecret, code)
	if err != nil {
		return nil, err
	}
	// generating random token
	token, err := utils.GenerateRandomString(32)
	if err != nil {
		return nil, err
	}
	// save token and claims for 5 minutes
	b, err := json.Marshal(claims)
	if err != nil {
		return nil, err
	}
	err = u.Repo.SetTTL("token:"+token, b, time.Minute*5)
	if err != nil {
		return nil, err
	}
	// generate the access token
	return &oauth2.Token{
		AccessToken: token,
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(time.Minute * 5),
	}, nil
}

// AccessResourcesWithToken returns the resources that the client has access to
// with the given token.
func (u *UseCase) AccessResourcesWithToken(token string) (map[string]interface{}, error) {
	// get the claims
	res := u.Repo.Get("token:" + token)
	if res == nil {
		return nil, fmt.Errorf("could not get token")
	}
	var claims utils.AuthCodeClaims
	err := json.Unmarshal(res, &claims)
	if err != nil {
		return nil, err
	}
	// get the resources
	scopes := strings.Split(claims.Scopes, ",")
	resources := make(map[string]interface{})
	for _, scope := range scopes {
		switch scope {
		case constants.READ_USER_SCOPE:
			user, err := u.GetUser(claims.UserID, true)
			if err != nil {
				return nil, err
			}
			resources["profile"] = user
		}
	}
	return resources, nil
}
