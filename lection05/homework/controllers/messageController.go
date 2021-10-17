package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/FMyb/tfs-go-hw/lection05/homework/domain"
	"github.com/FMyb/tfs-go-hw/lection05/homework/models"
	"github.com/FMyb/tfs-go-hw/lection05/homework/utils"
)

func SendToPublic(w http.ResponseWriter, r *http.Request) {
	d, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var m domain.Message
	err = json.Unmarshal(d, &m)
	if err != nil || m.Text == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m.UserID = r.Context().Value(utils.KeyUserID).(uint)
	message := models.SendToPublic(m)
	res, err := json.Marshal(message)
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

func SendToUser(w http.ResponseWriter, r *http.Request) {
	d, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var m domain.Message
	err = json.Unmarshal(d, &m)
	if err != nil || m.Text == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m.UserID = r.Context().Value(utils.KeyUserID).(uint)
	message := models.SendToUser(m)
	res, err := json.Marshal(message)
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

func GetMessagesFromPublic(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	messages := models.GetMessagesFromPublic()
	res, err := json.Marshal(messages)
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

func GetUserMessages(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	q := r.URL.Query()
	offset, err := strconv.Atoi(q.Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}
	length, err := strconv.Atoi(q.Get("length"))
	if err != nil || length < 0 {
		length = -1
	}
	uID := r.Context().Value(utils.KeyUserID).(uint)
	message := models.GetUserMessages(uID, offset, length)
	res, err := json.Marshal(message)
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
