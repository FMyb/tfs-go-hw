package tickers

import (
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/domain"
)

type MockwsTicker struct {
}

func NewMockwsTicker() *MockwsTicker {
	return &MockwsTicker{}
}

func (m MockwsTicker) NextTicker() domain.ResponseStatus {
	return domain.SuccessResponseTicker{
		Time:      123,
		Feed:      "ticker",
		ProductID: "productId",
		Suspended: false,
		MarkPrice: 100,
	}
}

func (m MockwsTicker) Start(_ []string) {}

func (m MockwsTicker) Stop() chan interface{} {
	ch := make(chan interface{})
	ch <- nil
	return ch
}
