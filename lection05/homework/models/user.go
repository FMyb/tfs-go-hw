package models

import (
	"errors"

	"github.com/FMyb/tfs-go-hw/lection05/homework/domain"
	uRepository "github.com/FMyb/tfs-go-hw/lection05/homework/repositories"
)

func Create(user domain.User) (*domain.User, error) {
	findUser := uRepository.FindUserByEmail(user.Email)
	if findUser != nil {
		return nil, errors.New("this user email is used")
	}
	return uRepository.CreateUser(user), nil
}

func Login(user domain.User) (*domain.User, error) {
	findUser := uRepository.FindUserByEmail(user.Email)
	if findUser == nil {
		return nil, errors.New("can't find user with email")
	}
	if findUser.Password != user.Password {
		return nil, errors.New("wrong password")
	}
	return findUser, nil
}
