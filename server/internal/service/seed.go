package service

import (
	"context"
	"errors"

	"github.com/iwwwanow/garden/server/internal/model"
	"github.com/iwwwanow/garden/server/internal/repo"
	"github.com/jackc/pgx/v5"
)

type SeedService interface {
	GetSeeds(ctx context.Context, userID int) ([]*model.Seed, error)
	Share(ctx context.Context, fromUserID, toUserID, flowerID, quantity int) error
}

type seedService struct {
	seeds         repo.SeedRepo
	users         repo.UserRepo
	notifications repo.NotificationRepo
}

func NewSeedService(s repo.SeedRepo, u repo.UserRepo, n repo.NotificationRepo) SeedService {
	return &seedService{seeds: s, users: u, notifications: n}
}

func (s *seedService) GetSeeds(ctx context.Context, userID int) ([]*model.Seed, error) {
	return s.seeds.ListByUserID(ctx, userID)
}

func (s *seedService) Share(ctx context.Context, fromUserID, toUserID, flowerID, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	// verify recipient exists
	if _, err := s.users.GetByID(ctx, toUserID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	// check sender has enough
	seed, err := s.seeds.GetByUserAndFlower(ctx, fromUserID, flowerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrInsufficientSeeds
		}
		return err
	}
	if seed.Quantity < quantity {
		return ErrInsufficientSeeds
	}
	// deduct from sender
	if err := s.seeds.Upsert(ctx, fromUserID, flowerID, -quantity); err != nil {
		return err
	}
	// add to recipient
	if err := s.seeds.Upsert(ctx, toUserID, flowerID, quantity); err != nil {
		return err
	}
	_ = s.notifications.Create(ctx, toUserID, "seeds_received", map[string]any{
		"from_user_id": fromUserID,
		"flower_id":    flowerID,
		"quantity":     quantity,
	})
	return nil
}
