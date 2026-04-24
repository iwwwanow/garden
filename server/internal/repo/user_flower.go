package repo

import (
	"context"

	"github.com/iwwwanow/garden/server/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserFlowerRepo interface {
	Create(ctx context.Context, userID, flowerID int) (*model.UserFlower, error)
	GetByID(ctx context.Context, id int) (*model.UserFlower, error)
	ListByUserID(ctx context.Context, userID int) ([]*model.UserFlower, error)
	CountActiveByUserID(ctx context.Context, userID int) (int, error)
	ListActive(ctx context.Context) ([]*model.UserFlower, error)
	IncrementDay(ctx context.Context, id int) error
	SetDried(ctx context.Context, id int) error
	SetLastWatered(ctx context.Context, id int) error
}

type userFlowerRepo struct{ pool *pgxpool.Pool }

func NewUserFlowerRepo(pool *pgxpool.Pool) UserFlowerRepo { return &userFlowerRepo{pool: pool} }

const ufColumns = `id, user_id, flower_id, day, last_watered_at, is_dried, created_at`

func scanUserFlower(row interface{ Scan(...any) error }) (*model.UserFlower, error) {
	uf := &model.UserFlower{}
	err := row.Scan(&uf.ID, &uf.UserID, &uf.FlowerID, &uf.Day, &uf.LastWateredAt, &uf.IsDried, &uf.CreatedAt)
	if err != nil {
		return nil, err
	}
	return uf, nil
}

func (r *userFlowerRepo) Create(ctx context.Context, userID, flowerID int) (*model.UserFlower, error) {
	return scanUserFlower(r.pool.QueryRow(ctx,
		`INSERT INTO user_flowers (user_id, flower_id) VALUES ($1, $2) RETURNING `+ufColumns,
		userID, flowerID,
	))
}

func (r *userFlowerRepo) GetByID(ctx context.Context, id int) (*model.UserFlower, error) {
	return scanUserFlower(r.pool.QueryRow(ctx,
		`SELECT `+ufColumns+` FROM user_flowers WHERE id = $1`, id,
	))
}

func (r *userFlowerRepo) ListByUserID(ctx context.Context, userID int) ([]*model.UserFlower, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+ufColumns+` FROM user_flowers WHERE user_id = $1 ORDER BY created_at`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.UserFlower
	for rows.Next() {
		uf, err := scanUserFlower(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, uf)
	}
	return list, rows.Err()
}

func (r *userFlowerRepo) CountActiveByUserID(ctx context.Context, userID int) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM user_flowers WHERE user_id = $1 AND is_dried = false`, userID,
	).Scan(&count)
	return count, err
}

func (r *userFlowerRepo) ListActive(ctx context.Context) ([]*model.UserFlower, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+ufColumns+` FROM user_flowers WHERE is_dried = false ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.UserFlower
	for rows.Next() {
		uf, err := scanUserFlower(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, uf)
	}
	return list, rows.Err()
}

func (r *userFlowerRepo) IncrementDay(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, `UPDATE user_flowers SET day = day + 1 WHERE id = $1`, id)
	return err
}

func (r *userFlowerRepo) SetDried(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, `UPDATE user_flowers SET is_dried = true WHERE id = $1`, id)
	return err
}

func (r *userFlowerRepo) SetLastWatered(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE user_flowers SET last_watered_at = CURRENT_DATE WHERE id = $1`, id,
	)
	return err
}
