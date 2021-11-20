package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/FMyb/tfs-go-hw/lection05/homework/data"
	"github.com/FMyb/tfs-go-hw/lection05/homework/domain"
	"github.com/FMyb/tfs-go-hw/lection05/homework/models"
	"github.com/FMyb/tfs-go-hw/lection05/homework/utils"
)

type Users struct {
	data.UserDB
}

func NewUsers(userDB data.UserDB) *Users {
	return &Users{UserDB: userDB}
}

func (us *Users) Login(w http.ResponseWriter, r *http.Request) {
	d, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var u domain.User
	err = json.Unmarshal(d, &u)
	if err != nil || u.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := models.Login(u, us)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token := utils.GetToken(*user)
	user.Token = token
	res, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (us *Users) UserRegister(w http.ResponseWriter, r *http.Request) {
	d, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var u domain.User
	err = json.Unmarshal(d, &u)
	if err != nil || u.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := models.Create(u, us)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token := utils.GetToken(*user)
	user.Token = token
	res, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
