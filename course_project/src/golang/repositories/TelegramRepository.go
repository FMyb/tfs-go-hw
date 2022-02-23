package repositories

import (
	"context"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/repositories/queries"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ClientRepository interface {
	SaveUser(ctx context.Context, userID int64) error
	DeleteUser(ctx context.Context, userID int64) error
	Users(ctx context.Context) ([]int64, error)
}

type telegramRepo struct {
	*queries.Queries
	pool *pgxpool.Pool
}

func NewTelegramRepository(pgxPool *pgxpool.Pool) ClientRepository {
	return &telegramRepo{
		Queries: queries.New(pgxPool),
		pool:    pgxPool,
	}
}
