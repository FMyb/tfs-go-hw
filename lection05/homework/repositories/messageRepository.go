package repositories

import (
	"github.com/FMyb/tfs-go-hw/lection05/homework/data"
	"github.com/FMyb/tfs-go-hw/lection05/homework/domain"
)

func SendToPublic(message domain.Message, mes data.MessagesDB) *domain.Message {
	message.ID = mes.MessageIDGetAndInc()
	message.ToUserID = 0
	mes.AddMessage(message)
	return &message
}

func SendToUser(message domain.Message, mes data.MessagesDB) *domain.Message {
	message.ID = mes.MessageIDGetAndInc()
	if message.ToUserID == 0 {
		return nil
	}
	mes.AddMessage(message)
	return &message
}

func GetMessagesFromPublic(mes data.MessagesDB) []domain.Message {
	result := make([]domain.Message, 0)
	for _, message := range mes.Messages() {
		if message.ToUserID == 0 {
			result = append(result, message)
		}
	}
	return result
}

func GetUserMessages(userID uint, offset int, length int, mes data.MessagesDB) []domain.Message {
	result := make([]domain.Message, 0)
	cnt := 0
	for _, message := range mes.Messages() {
		if length != -1 && cnt > length {
			break
		}
		if cnt < offset {
			cnt++
			continue
		}
		if message.ToUserID == userID || message.UserID == userID {
			result = append(result, message)
			cnt++
		}
	}
	return result
}
