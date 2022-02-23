package services

import (
	"context"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/domain"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/repositories"
)

type MockBot struct {
	tickers    Tickers
	repository repositories.OrderRepository
	ctx        context.Context
	client     Client
	orders     Orders
	cancel     context.CancelFunc
}

func NewMockBot(tickers Tickers, repository repositories.OrderRepository, ctx context.Context, client Client, orders Orders, cancel context.CancelFunc) *MockBot {
	return &MockBot{tickers: tickers, repository: repository, ctx: ctx, client: client, orders: orders, cancel: cancel}
}

func (m MockBot) Orders() Orders {
	return m.orders
}

func (m MockBot) Tickers() Tickers {
	return m.tickers
}

func (m MockBot) ProductId() string {
	return "productId"
}

func (m MockBot) PublicApiKey() string {
	return "publicApiKey"
}

func (m MockBot) Repository() repositories.OrderRepository {
	return m.repository
}

func (m MockBot) Context() context.Context {
	return m.ctx
}

func (m MockBot) Client() Client {
	return m.client
}

func (m MockBot) PrivateApiKey() string {
	return "privateApiKey"
}

func (m MockBot) Sell(_ domain.SuccessResponseTicker) bool {
	return true
}

func (m MockBot) Buy(_ domain.SuccessResponseTicker) bool {
	return true
}

func (m MockBot) SetLastPrice(_ float32) {}

func (m MockBot) Cancel() context.CancelFunc {
	return m.cancel
}
