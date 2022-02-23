package orders

import (
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/domain"
	"time"
)

type MockOrder struct{}

func NewMockOrder() *MockOrder {
	return &MockOrder{}
}

func (m MockOrder) SendOrder(productId string, size uint64, side string,
	publicApiKey string, privateApiKey string) (domain.ResponseOrder, error) {
	return domain.KrakenResponseOrder{
		Kresult:     domain.SUCCESSES,
		Kerror:      "",
		KserverTime: time.Now(),
		SendStatus: domain.SendStatus{
			OrderId: "orderId",
			Status:  "placed",
			OrderEvents: []domain.OrderEvent{
				{
					ExecutionOrderEvent: domain.ExecutionOrderEvent{
						ExecPrice:  100,
						ExecAmount: 1,
						OrderPriorExecution: struct {
							Symbol string `json:"symbol"`
							Side   string `json:"side"`
						}{
							Symbol: "symbol",
							Side:   "buy",
						},
					},
					OType: domain.EXECUTION,
				},
			},
		},
	}, nil
}
