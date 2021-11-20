package models

import (
	"errors"

	"github.com/FMyb/tfs-go-hw/lection05/homework/data"
	"github.com/FMyb/tfs-go-hw/lection05/homework/domain"
	uRepository "github.com/FMyb/tfs-go-hw/lection05/homework/repositories"
)

func Create(user domain.User, us data.UserDB) (*domain.User, error) {
	findUser := uRepository.FindUserByEmail(user.Email, us)
	if findUser != nil {
		return nil, errors.New("this user email is used")
	}
	return uRepository.CreateUser(user, us), nil
}

func Login(user domain.User, us data.UserDB) (*domain.User, error) {
	findUser := uRepository.FindUserByEmail(user.Email, us)
	if findUser == nil {
		return nil, errors.New("can't find user with email")
	}
	if findUser.Password != user.Password {
		return nil, errors.New("wrong password")
	}
	return findUser, nil
}
