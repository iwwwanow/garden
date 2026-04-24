package service

import (
	"context"
	"log"
	"time"

	"github.com/iwwwanow/garden/server/internal/repo"
)

type TickService interface {
	// RunTick processes waterings for checkDate and advances all active flowers.
	// Real scheduler calls with yesterday; dev endpoint calls with today.
	RunTick(ctx context.Context, checkDate time.Time) error
}

type tickService struct {
	userFlowers   repo.UserFlowerRepo
	waterings     repo.WateringRepo
	users         repo.UserRepo
	seeds         repo.SeedRepo
	notifications repo.NotificationRepo
}

func NewTickService(
	uf repo.UserFlowerRepo,
	w repo.WateringRepo,
	u repo.UserRepo,
	s repo.SeedRepo,
	n repo.NotificationRepo,
) TickService {
	return &tickService{userFlowers: uf, waterings: w, users: u, seeds: s, notifications: n}
}

func (s *tickService) RunTick(ctx context.Context, checkDate time.Time) error {
	flowers, err := s.userFlowers.ListActive(ctx)
	if err != nil {
		return err
	}
	log.Printf("tick: processing %d active flowers for date %s", len(flowers), checkDate.Format("2006-01-02"))

	for _, uf := range flowers {
		watered, err := s.waterings.AnyOnDate(ctx, uf.ID, checkDate)
		if err != nil {
			log.Printf("tick: check watering flower=%d: %v", uf.ID, err)
			continue
		}

		if !watered {
			if err := s.userFlowers.SetDried(ctx, uf.ID); err != nil {
				log.Printf("tick: set dried flower=%d: %v", uf.ID, err)
				continue
			}
			_ = s.notifications.Create(ctx, uf.UserID, "flower_dried", map[string]any{
				"user_flower_id": uf.ID,
				"flower_id":      uf.FlowerID,
			})
			continue
		}

		// Flower survived: award FD = current day value
		if err := s.users.AddFD(ctx, uf.UserID, uf.Day); err != nil {
			log.Printf("tick: add fd user=%d: %v", uf.UserID, err)
		}

		// Advance day
		if err := s.userFlowers.IncrementDay(ctx, uf.ID); err != nil {
			log.Printf("tick: increment day flower=%d: %v", uf.ID, err)
			continue
		}

		newDay := uf.Day + 1
		if newDay%7 == 0 {
			if err := s.seeds.Upsert(ctx, uf.UserID, uf.FlowerID, 1); err != nil {
				log.Printf("tick: give seed user=%d flower=%d: %v", uf.UserID, uf.FlowerID, err)
			}
			_ = s.notifications.Create(ctx, uf.UserID, "seed_earned", map[string]any{
				"user_flower_id": uf.ID,
				"flower_id":      uf.FlowerID,
				"day":            newDay,
			})
		}
	}
	return nil
}
