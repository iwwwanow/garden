package repo

import (
	"context"
	"encoding/json"

	"github.com/iwwwanow/garden/server/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepo interface {
	Create(ctx context.Context, userID int, notifType string, payload map[string]any) error
	ListByUserID(ctx context.Context, userID int) ([]*model.Notification, error)
}

type notificationRepo struct{ pool *pgxpool.Pool }

func NewNotificationRepo(pool *pgxpool.Pool) NotificationRepo {
	return &notificationRepo{pool: pool}
}

func (r *notificationRepo) Create(ctx context.Context, userID int, notifType string, payload map[string]any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx,
		`INSERT INTO notifications (user_id, type, payload) VALUES ($1, $2, $3)`,
		userID, notifType, data,
	)
	return err
}

func (r *notificationRepo) ListByUserID(ctx context.Context, userID int) ([]*model.Notification, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, type, payload, is_read, created_at
		 FROM notifications WHERE user_id = $1 ORDER BY created_at DESC LIMIT 100`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.Notification
	for rows.Next() {
		n := &model.Notification{}
		var payload []byte
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &payload, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, err
		}
		n.Payload = json.RawMessage(payload)
		list = append(list, n)
	}
	return list, rows.Err()
}
