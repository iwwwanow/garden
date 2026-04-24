package repo

import (
	"context"

	"github.com/iwwwanow/garden/server/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FlowerRepo interface {
	GetByID(ctx context.Context, id int) (*model.Flower, error)
	ListAll(ctx context.Context) ([]*model.Flower, error)
}

type flowerRepo struct{ pool *pgxpool.Pool }

func NewFlowerRepo(pool *pgxpool.Pool) FlowerRepo { return &flowerRepo{pool: pool} }

func scanFlower(row interface{ Scan(...any) error }) (*model.Flower, error) {
	f := &model.Flower{}
	err := row.Scan(&f.ID, &f.Season, &f.ImagePath, &f.CreatedAt)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (r *flowerRepo) GetByID(ctx context.Context, id int) (*model.Flower, error) {
	return scanFlower(r.pool.QueryRow(ctx,
		`SELECT id, season, image_path, created_at FROM flowers WHERE id = $1`, id,
	))
}

func (r *flowerRepo) ListAll(ctx context.Context) ([]*model.Flower, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, season, image_path, created_at FROM flowers ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var flowers []*model.Flower
	for rows.Next() {
		f, err := scanFlower(rows)
		if err != nil {
			return nil, err
		}
		flowers = append(flowers, f)
	}
	return flowers, rows.Err()
}
