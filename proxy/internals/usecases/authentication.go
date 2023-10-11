package usecases

import (
	"encoding/json"
	"fmt"
	utilities "github.com/globe-and-citizen/layer8-utils"
	"layer8-proxy/config"
	"layer8-proxy/constants"
	"layer8-proxy/entities"

	"golang.org/x/crypto/bcrypt"
)

// GetUserByToken returns the user associated with the given token
func (u *UseCase) GetUserByToken(token string) (*entities.User, error) {
	// verify token
	userID, err := utilities.VerifyUserToken(config.SECRET_KEY, token)
	if err != nil {
		return nil, err
	}
	// get user
	user, err := u.GetUser(userID, false)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// LoginUser logs in a user and returns a map containing the access token and the user data
func (u *UseCase) LoginUser(username, password string) (map[string]interface{}, error) {
	// find user
	users, err := u.Repo.Keys("^user:*")
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, constants.ErrNotFound
	}

	var user *entities.User
	for i, v := range users {
		buser := u.Repo.Get(v)
		err := json.Unmarshal(buser, &user)
		if err != nil {
			return nil, err
		}
		if user.Username == username {
			break
		}
		if i == len(users)-1 {
			return nil, constants.ErrNotFound
		}
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, constants.ErrInvalidPassword
	}
	// generate token
	token, err := utilities.GenerateUserToken(config.SECRET_KEY, user.ID)
	if err != nil {
		return nil, err
	}
	// return token and user
	return map[string]interface{}{
		"token": token,
		"user":  user,
	}, nil
}

// RegisterUser registers a user and returns a map containing the access token and the user data
func (u *UseCase) RegisterUser(user *entities.User) (map[string]interface{}, error) {
	// check if user already exists
	users, err := u.Repo.Keys("^user:*")
	if err != nil {
		return nil, err
	}
	for _, v := range users {
		var u2 entities.User
		res := u.Repo.Get(v)
		err := json.Unmarshal(res, &u2)
		if err != nil {
			return nil, err
		}
		if u2.Email == user.Email {
			return nil, constants.ErrAlreadyExists
		}
	}
	// encrypt password
	bypass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	// Generate new ID
	id, err := u.Repo.Incr("_user:id")
	if err != nil {
		return nil, err
	}
	user.ID = id
	user.Password = string(bypass)
	// create pseudonymized data
	pname, pmail, pfname, plname := utilities.GeneratePlaceholderUserData()
	user.PsedonymizedData = entities.AbstractUser{
		Username: pname,
		Email:    pmail,
		Fname:    pfname,
		Lname:    plname,
	}
	// save user
	buser, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	err = u.Repo.Set(fmt.Sprintf("user:%d", user.ID), buser)
	if err != nil {
		return nil, err
	}
	// generate token
	token, err := utilities.GenerateUserToken(config.SECRET_KEY, user.ID)
	if err != nil {
		return nil, err
	}
	// return token and user
	return map[string]interface{}{
		"token": token,
		"user":  user,
	}, nil
}
