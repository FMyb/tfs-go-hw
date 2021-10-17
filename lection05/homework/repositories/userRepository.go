package repositories

import (
	"github.com/FMyb/tfs-go-hw/lection05/homework/domain"
)

var users []domain.User
var userID uint = 1

func FindUserByEmail(email string) *domain.User {
	for _, user := range users {
		if user.Email == email {
			return &user
		}
	}
	return nil
}

func CreateUser(user domain.User) *domain.User {
	user.ID = userID
	userID++
	users = append(users, user)
	return &user
}
