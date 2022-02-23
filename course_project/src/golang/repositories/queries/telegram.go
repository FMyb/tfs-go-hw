package queries

import (
	"context"
	log "github.com/sirupsen/logrus"
)

const saveUserQuery = `INSERT INTO client_users (user_id) VALUES ($1);`

func (q *Queries) SaveUser(ctx context.Context, chatID int64) error {
	_, err := q.pool.Query(
		ctx,
		saveUserQuery,
		chatID,
	)
	if err != nil {
		return err
	}
	return nil
}

const deleteUserQuery = `DELETE FROM client_users WHERE user_id = $1;`

func (q *Queries) DeleteUser(ctx context.Context, chatID int64) error {
	_, err := q.pool.Query(
		ctx,
		deleteUserQuery,
		chatID,
	)
	if err != nil {
		return err
	}
	return nil
}

const usersQuery = `SELECT user_id FROM client_users`

func (q *Queries) Users(ctx context.Context) ([]int64, error) {
	log.Debug("select all users")
	rows, err := q.pool.Query(
		ctx,
		usersQuery,
	)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	chatIDs := make([]int64, 0)
	log.Debug("start scan users")
	for rows.Next() {
		var chatID int
		err = rows.Scan(&chatID)
		log.Debugf("scan user: %d", chatID)
		if err != nil {
			return nil, err
		}
		chatIDs = append(chatIDs, int64(chatID))
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return chatIDs, nil
}
