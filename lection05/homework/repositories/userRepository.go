package repositories

import (
	"github.com/FMyb/tfs-go-hw/lection05/homework/data"
	"github.com/FMyb/tfs-go-hw/lection05/homework/domain"
)

func FindUserByEmail(email string, us data.UserDB) *domain.User {
	for _, user := range us.Users() {
		if user.Email == email {
			return &user
		}
	}
	return nil
}

func CreateUser(user domain.User, us data.UserDB) *domain.User {
	user.ID = us.UsersIDGetAndInc()
	us.AddUser(user)
	return &user
}
