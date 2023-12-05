package usecases

import (
	"fmt"
	"globe-and-citizen/layer8/proxy/config"
	"globe-and-citizen/layer8/proxy/models"
	"globe-and-citizen/layer8/proxy/utils"

	utilities "github.com/globe-and-citizen/layer8-utils"
)

// GetUserByToken returns the user associated with the given token
func (u *UseCase) GetUserByToken(token string) (*models.User, error) {
	// verify token
	userID, err := utilities.VerifyUserToken(config.SECRET_KEY, token)
	if err != nil {
		return nil, err
	}
	// get user
	user, err := u.Repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UseCase) LoginUser(username, password string) (map[string]interface{}, error) {

	userSalt, err := u.Repo.LoginUserPrecheck(username)
	if err != nil {
		return nil, err
	}

	HashedAndSaltedPass := utils.SaltAndHashPassword(password, userSalt)

	user, err := u.Repo.GetUser(username)
	if err != nil {
		return nil, err
	}

	if user.Password != HashedAndSaltedPass {
		return nil, fmt.Errorf("invalid password")
	}

	token, err := utilities.GenerateUserToken(config.SECRET_KEY, int64(user.ID))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token": token,
		"user":  user,
	}, nil
}
