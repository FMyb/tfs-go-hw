package data

import (
	"sync"

	"github.com/FMyb/tfs-go-hw/lection05/homework/domain"
)

type UserDB interface {
	AddUser(user domain.User)
	Users() []domain.User
	UsersIDGetAndInc() uint
	UserID() uint
	SetUserID(userID uint)
}

type Users struct {
	users  []domain.User
	userID uint
	sync.Mutex
}

func NewUsers() *Users {
	return &Users{
		users:  make([]domain.User, 0),
		userID: 1,
		Mutex:  sync.Mutex{},
	}
}

func (us *Users) AddUser(user domain.User) {
	us.Lock()
	us.users = append(us.users, user)
	us.Unlock()
}

func (us *Users) Users() []domain.User {
	us.Lock()
	defer us.Unlock()
	return us.users
}

func (us *Users) UsersIDGetAndInc() uint {
	us.Lock()
	defer us.Unlock()
	defer func() {
		us.userID++
	}()
	return us.userID
}

func (us *Users) UserID() uint {
	us.Lock()
	defer us.Unlock()
	return us.userID
}

func (us *Users) SetUserID(userID uint) {
	us.Lock()
	us.userID = userID
	us.Unlock()
}
