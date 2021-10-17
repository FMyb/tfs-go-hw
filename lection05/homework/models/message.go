package models

import (
	"github.com/FMyb/tfs-go-hw/lection05/homework/domain"
	"github.com/FMyb/tfs-go-hw/lection05/homework/repositories"
)

func SendToPublic(message domain.Message) *domain.Message {
	return repositories.SendToPublic(message)
}

func SendToUser(message domain.Message) *domain.Message {
	return repositories.SendToUser(message)
}

func GetMessagesFromPublic() []domain.Message {
	return repositories.GetMessagesFromPublic()
}

func GetUserMessages(userID uint, offset int, length int) []domain.Message {
	return repositories.GetUserMessages(userID, offset, length)
}
