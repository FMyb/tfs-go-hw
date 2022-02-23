package repositories

import (
	"context"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/domain"
)

type MockOrderRepository struct{}

func (m MockOrderRepository) SendOrder(ctx context.Context, order domain.ResponseOrder) error {
	return nil
}

func NewMockOrderRepository() *MockOrderRepository {
	return &MockOrderRepository{}
}
