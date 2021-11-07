package data

import (
	"sync"

	"github.com/FMyb/tfs-go-hw/lection05/homework/domain"
)

type MessagesDB interface {
	AddMessage(message domain.Message)
	Messages() []domain.Message
	MessageIDGetAndInc() uint
	MessageID() uint
	SetMessageID(messageID uint)
}

type Messages struct {
	messages  []domain.Message
	messageID uint
	sync.Mutex
}

func NewMessages() *Messages {
	return &Messages{
		messages:  make([]domain.Message, 0),
		messageID: 1,
		Mutex:     sync.Mutex{},
	}
}

func (mes *Messages) AddMessage(message domain.Message) {
	mes.Lock()
	mes.messages = append(mes.messages, message)
	mes.Unlock()
}

func (mes *Messages) Messages() []domain.Message {
	mes.Lock()
	defer mes.Unlock()
	return mes.messages
}

func (mes *Messages) MessageIDGetAndInc() uint {
	mes.Lock()
	defer mes.Unlock()
	defer func() {
		mes.messageID++
	}()
	return mes.messageID
}

func (mes *Messages) MessageID() uint {
	mes.Lock()
	defer mes.Unlock()
	return mes.messageID
}

func (mes *Messages) SetMessageID(messageID uint) {
	mes.Lock()
	mes.messageID = messageID
	mes.Unlock()
}
