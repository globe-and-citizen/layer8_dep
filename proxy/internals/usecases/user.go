package usecases

import (
	"encoding/json"
	"fmt"
	"globe-and-citizen/layer8/proxy/constants"
	"globe-and-citizen/layer8/proxy/entities"

	"golang.org/x/crypto/bcrypt"
)

func (u *UseCase) GetUser(id int64, pseudonymized bool) (*entities.User, error) {
	var user entities.User
	res := u.Repo.Get(fmt.Sprintf("user:%d", id))
	if res == nil {
		return nil, constants.ErrNotFound
	}
	err := json.Unmarshal(res, &user)
	if err != nil {
		return nil, err
	}
	if !pseudonymized {
		return &user, nil
	}
	return &entities.User{
		ID: user.ID,
		AbstractUser: entities.AbstractUser{
			Username: user.PsedonymizedData.Username,
			Email:    user.PsedonymizedData.Email,
			Fname:    user.PsedonymizedData.Fname,
			Lname:    user.PsedonymizedData.Lname,
		},
	}, nil
}

func (u *UseCase) AddUser(user *entities.User) (*entities.User, error) {
	// check if user already exists
	users, err := u.Repo.Keys("user:*")
	if err != nil {
		return nil, err
	}
	for _, v := range users {
		var u2 entities.User
		res := u.Repo.Get(v)
		err := json.Unmarshal(res, &u)
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
	// Previously create pseudonymized data
	// Now, that pseudoanonymized data hase been replaced by the variable names
	user.PsedonymizedData.Username = "pname"
	user.PsedonymizedData.Email = "pmail"
	user.PsedonymizedData.Fname = "pfname"
	user.PsedonymizedData.Lname = "plname"
	// save user
	b, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	err = u.Repo.Set(fmt.Sprintf("user:%d", user.ID), b)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UseCase) DeleteUser(id int64) error {
	return u.Repo.Delete(fmt.Sprintf("user:%d", id))
}
