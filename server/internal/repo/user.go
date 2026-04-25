package repo

import (
	"context"

	"github.com/iwwwanow/garden/server/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo interface {
	Create(ctx context.Context, username, passwordHash, firstName string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByID(ctx context.Context, id int) (*model.User, error)
	AddFD(ctx context.Context, userID, delta int) error
	UpdateFirstName(ctx context.Context, userID int, firstName string) error
	Leaderboard(ctx context.Context, limit, offset int) ([]*model.User, error)
	ListAll(ctx context.Context) ([]*model.User, error)
}

type userRepo struct{ pool *pgxpool.Pool }

func NewUserRepo(pool *pgxpool.Pool) UserRepo { return &userRepo{pool: pool} }

const userColumns = `id, username, password_hash, first_name, telegram_id, fd_balance, created_at`

func scanUser(row interface{ Scan(...any) error }) (*model.User, error) {
	u := &model.User{}
	err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.FirstName, &u.TelegramID, &u.FDBalance, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepo) Create(ctx context.Context, username, passwordHash, firstName string) (*model.User, error) {
	row := r.pool.QueryRow(ctx,
		`INSERT INTO users (username, password_hash, first_name) VALUES ($1, $2, $3) RETURNING `+userColumns,
		username, passwordHash, firstName,
	)
	return scanUser(row)
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT `+userColumns+` FROM users WHERE username = $1`, username,
	)
	return scanUser(row)
}

func (r *userRepo) GetByID(ctx context.Context, id int) (*model.User, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT `+userColumns+` FROM users WHERE id = $1`, id,
	)
	return scanUser(row)
}

func (r *userRepo) AddFD(ctx context.Context, userID, delta int) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET fd_balance = fd_balance + $1 WHERE id = $2`, delta, userID,
	)
	return err
}

func (r *userRepo) UpdateFirstName(ctx context.Context, userID int, firstName string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET first_name = $1 WHERE id = $2`, firstName, userID,
	)
	return err
}

func (r *userRepo) Leaderboard(ctx context.Context, limit, offset int) ([]*model.User, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+userColumns+` FROM users ORDER BY fd_balance DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*model.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *userRepo) ListAll(ctx context.Context) ([]*model.User, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+userColumns+` FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*model.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}
