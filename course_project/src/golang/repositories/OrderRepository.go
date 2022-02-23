package repositories

import (
	"context"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/domain"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/repositories/queries"
	"github.com/jackc/pgx/v4/pgxpool"
)

type OrderRepository interface {
	SendOrder(ctx context.Context, order domain.ResponseOrder) error
}

type orderRepo struct {
	*queries.Queries
	pool *pgxpool.Pool
}

func NewOrderRepository(pgxPool *pgxpool.Pool) OrderRepository {
	return &orderRepo{
		Queries: queries.New(pgxPool),
		pool:    pgxPool,
	}
}
