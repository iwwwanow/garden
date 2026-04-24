package service

import (
	"context"
	"errors"
	"time"

	"github.com/iwwwanow/garden/server/internal/model"
	"github.com/iwwwanow/garden/server/internal/repo"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const maxActiveFlowers = 64

type FlowerService interface {
	GetUserFlowers(ctx context.Context, userID int) ([]*model.UserFlower, error)
	Water(ctx context.Context, userFlowerID, wateredByUserID int) error
	Plant(ctx context.Context, userID, flowerID int) (*model.UserFlower, error)
}

type flowerService struct {
	userFlowers   repo.UserFlowerRepo
	waterings     repo.WateringRepo
	seeds         repo.SeedRepo
	flowers       repo.FlowerRepo
	notifications repo.NotificationRepo
}

func NewFlowerService(
	uf repo.UserFlowerRepo,
	w repo.WateringRepo,
	s repo.SeedRepo,
	f repo.FlowerRepo,
	n repo.NotificationRepo,
) FlowerService {
	return &flowerService{userFlowers: uf, waterings: w, seeds: s, flowers: f, notifications: n}
}

func (s *flowerService) GetUserFlowers(ctx context.Context, userID int) ([]*model.UserFlower, error) {
	return s.userFlowers.ListByUserID(ctx, userID)
}

func (s *flowerService) Water(ctx context.Context, userFlowerID, wateredByUserID int) error {
	uf, err := s.userFlowers.GetByID(ctx, userFlowerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	if uf.IsDried {
		return ErrFlowerDried
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	already, err := s.waterings.HasWatered(ctx, userFlowerID, wateredByUserID, today)
	if err != nil {
		return err
	}
	if already {
		return ErrAlreadyWatered
	}

	if err := s.waterings.Record(ctx, userFlowerID, wateredByUserID, today); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrAlreadyWatered
		}
		return err
	}

	if err := s.userFlowers.SetLastWatered(ctx, userFlowerID); err != nil {
		return err
	}

	// notify owner if watered by someone else
	if uf.UserID != wateredByUserID {
		_ = s.notifications.Create(ctx, uf.UserID, "flower_watered", map[string]any{
			"user_flower_id":    userFlowerID,
			"watered_by_user_id": wateredByUserID,
		})
	}
	return nil
}

func (s *flowerService) Plant(ctx context.Context, userID, flowerID int) (*model.UserFlower, error) {
	// verify flower template exists
	if _, err := s.flowers.GetByID(ctx, flowerID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// check active flower limit
	count, err := s.userFlowers.CountActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if count >= maxActiveFlowers {
		return nil, ErrMaxFlowers
	}

	// spend 1 seed
	seed, err := s.seeds.GetByUserAndFlower(ctx, userID, flowerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInsufficientSeeds
		}
		return nil, err
	}
	if seed.Quantity < 1 {
		return nil, ErrInsufficientSeeds
	}
	if err := s.seeds.Upsert(ctx, userID, flowerID, -1); err != nil {
		return nil, err
	}

	return s.userFlowers.Create(ctx, userID, flowerID)
}
