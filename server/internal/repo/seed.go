package repo

import (
	"context"

	"github.com/iwwwanow/garden/server/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SeedRepo interface {
	ListByUserID(ctx context.Context, userID int) ([]*model.Seed, error)
	GetByUserAndFlower(ctx context.Context, userID, flowerID int) (*model.Seed, error)
	Upsert(ctx context.Context, userID, flowerID, delta int) error
}

type seedRepo struct{ pool *pgxpool.Pool }

func NewSeedRepo(pool *pgxpool.Pool) SeedRepo { return &seedRepo{pool: pool} }

func (r *seedRepo) ListByUserID(ctx context.Context, userID int) ([]*model.Seed, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, flower_id, quantity FROM seeds WHERE user_id = $1 ORDER BY flower_id`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var seeds []*model.Seed
	for rows.Next() {
		s := &model.Seed{}
		if err := rows.Scan(&s.ID, &s.UserID, &s.FlowerID, &s.Quantity); err != nil {
			return nil, err
		}
		seeds = append(seeds, s)
	}
	return seeds, rows.Err()
}

func (r *seedRepo) GetByUserAndFlower(ctx context.Context, userID, flowerID int) (*model.Seed, error) {
	s := &model.Seed{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, flower_id, quantity FROM seeds WHERE user_id=$1 AND flower_id=$2`,
		userID, flowerID,
	).Scan(&s.ID, &s.UserID, &s.FlowerID, &s.Quantity)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Upsert adds delta to the seed quantity (delta can be negative for spending).
func (r *seedRepo) Upsert(ctx context.Context, userID, flowerID, delta int) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO seeds (user_id, flower_id, quantity)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (user_id, flower_id)
		 DO UPDATE SET quantity = seeds.quantity + $3`,
		userID, flowerID, delta,
	)
	return err
}
