package models

import (
	"github.com/FMyb/tfs-go-hw/lection05/homework/data"
	"github.com/FMyb/tfs-go-hw/lection05/homework/domain"
	"github.com/FMyb/tfs-go-hw/lection05/homework/repositories"
)

func SendToPublic(message domain.Message, mes data.MessagesDB) *domain.Message {
	return repositories.SendToPublic(message, mes)
}

func SendToUser(message domain.Message, mes data.MessagesDB) *domain.Message {
	return repositories.SendToUser(message, mes)
}

func GetMessagesFromPublic(mes data.MessagesDB) []domain.Message {
	return repositories.GetMessagesFromPublic(mes)
}

func GetUserMessages(userID uint, offset int, length int, mes data.MessagesDB) []domain.Message {
	return repositories.GetUserMessages(userID, offset, length, mes)
}
