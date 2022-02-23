package telegram

import (
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/domain"
)

type MockClient struct{}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (m MockClient) SendOrder(order domain.ResponseOrder) error {
	return nil
}
