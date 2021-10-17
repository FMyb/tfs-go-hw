package repositories

import "github.com/FMyb/tfs-go-hw/lection05/homework/domain"

var messages []domain.Message

var messageID uint = 1

func SendToPublic(message domain.Message) *domain.Message {
	message.ID = messageID
	messageID++
	message.ToUserID = 0
	messages = append(messages, message)
	return &message
}

func SendToUser(message domain.Message) *domain.Message {
	message.ID = messageID
	messageID++
	if message.ToUserID == 0 {
		return nil
	}
	messages = append(messages, message)
	return &message
}

func GetMessagesFromPublic() []domain.Message {
	result := make([]domain.Message, 0)
	for _, message := range messages {
		if message.ToUserID == 0 {
			result = append(result, message)
		}
	}
	return result
}

func GetUserMessages(userID uint, offset int, length int) []domain.Message {
	result := make([]domain.Message, 0)
	cnt := 0
	for _, message := range messages {
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
