package repo

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WateringRepo interface {
	Record(ctx context.Context, userFlowerID, wateredByUserID int, date time.Time) error
	HasWatered(ctx context.Context, userFlowerID, wateredByUserID int, date time.Time) (bool, error)
	AnyOnDate(ctx context.Context, userFlowerID int, date time.Time) (bool, error)
}

type wateringRepo struct{ pool *pgxpool.Pool }

func NewWateringRepo(pool *pgxpool.Pool) WateringRepo { return &wateringRepo{pool: pool} }

func (r *wateringRepo) Record(ctx context.Context, userFlowerID, wateredByUserID int, date time.Time) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO waterings (user_flower_id, watered_by_user_id, watered_date) VALUES ($1, $2, $3)`,
		userFlowerID, wateredByUserID, date,
	)
	return err
}

func (r *wateringRepo) HasWatered(ctx context.Context, userFlowerID, wateredByUserID int, date time.Time) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM waterings WHERE user_flower_id=$1 AND watered_by_user_id=$2 AND watered_date=$3)`,
		userFlowerID, wateredByUserID, date,
	).Scan(&exists)
	return exists, err
}

func (r *wateringRepo) AnyOnDate(ctx context.Context, userFlowerID int, date time.Time) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM waterings WHERE user_flower_id=$1 AND watered_date=$2)`,
		userFlowerID, date,
	).Scan(&exists)
	return exists, err
}
