package queries

import (
	"context"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/domain"
)

const sendOrderQuery = `INSERT INTO sended_order (result, order_id, side, status, symbol, quantity, price, type, server_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`

func (q *Queries) SendOrder(ctx context.Context, order domain.ResponseOrder) error {
	rows, err := q.pool.Query(
		ctx,
		sendOrderQuery,
		order.Result(),
		order.OrderId(),
		order.Side(),
		order.Status(),
		order.Symbol(),
		order.Quantity(),
		order.Price(),
		order.Type(),
		order.ServerTime(),
	)
	defer rows.Close()
	if err != nil {
		return err // TODO добавить обработку
	}
	if rows.Err() != nil {
		return rows.Err()
	}
	return nil
}
