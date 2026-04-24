package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/iwwwanow/garden/server/internal/repo/mock"
	"github.com/iwwwanow/garden/server/internal/service"
)

func newTickSvc(
	uf *mock.UserFlowerRepo,
	w *mock.WateringRepo,
	u *mock.UserRepo,
	s *mock.SeedRepo,
	n *mock.NotificationRepo,
) service.TickService {
	return service.NewTickService(uf, w, u, s, n)
}

func TestTick_FlowerDriesWhenNotWatered(t *testing.T) {
	uf := mock.NewUserFlowerRepo()
	w := mock.NewWateringRepo()
	u := mock.NewUserRepo()
	s := mock.NewSeedRepo()
	n := mock.NewNotificationRepo()

	u.Create(context.Background(), "alice", "hash", "Alice")
	uf.Create(context.Background(), 1, 1)

	svc := newTickSvc(uf, w, u, s, n)
	// no watering recorded — flower should die
	if err := svc.RunTick(context.Background(), todayDate()); err != nil {
		t.Fatalf("tick error: %v", err)
	}

	flower := uf.ByID(1)
	if !flower.IsDried {
		t.Error("expected flower to be dried")
	}
	if n.OfType(1, "flower_dried") == nil {
		t.Error("expected flower_dried notification")
	}
}

func TestTick_FDAndDayIncrementWhenWatered(t *testing.T) {
	uf := mock.NewUserFlowerRepo()
	w := mock.NewWateringRepo()
	u := mock.NewUserRepo()
	s := mock.NewSeedRepo()
	n := mock.NewNotificationRepo()

	u.Create(context.Background(), "alice", "hash", "Alice")
	uf.Create(context.Background(), 1, 1) // flower is on day 1

	checkDate := todayDate()
	w.Record(context.Background(), 1, 1, checkDate) // mark as watered

	svc := newTickSvc(uf, w, u, s, n)
	if err := svc.RunTick(context.Background(), checkDate); err != nil {
		t.Fatalf("tick error: %v", err)
	}

	flower := uf.ByID(1)
	if flower.IsDried {
		t.Error("flower should not be dried")
	}
	if flower.Day != 2 {
		t.Errorf("day: got %d, want 2", flower.Day)
	}

	user := u.ByID(1)
	if user.FDBalance != 1 { // FD(day=1) = 1
		t.Errorf("fd_balance: got %d, want 1", user.FDBalance)
	}
}

func TestTick_SeedGivenOnDay7(t *testing.T) {
	uf := mock.NewUserFlowerRepo()
	w := mock.NewWateringRepo()
	u := mock.NewUserRepo()
	s := mock.NewSeedRepo()
	n := mock.NewNotificationRepo()

	u.Create(context.Background(), "alice", "hash", "Alice")
	flower, _ := uf.Create(context.Background(), 1, 1)

	// Manually set flower to day 6 (tick will bring it to 7)
	for i := 0; i < 5; i++ {
		uf.IncrementDay(context.Background(), flower.ID)
	}

	checkDate := todayDate()
	w.Record(context.Background(), flower.ID, 1, checkDate)

	svc := newTickSvc(uf, w, u, s, n)
	if err := svc.RunTick(context.Background(), checkDate); err != nil {
		t.Fatalf("tick error: %v", err)
	}

	if s.Quantity(1, 1) != 1 {
		t.Errorf("expected 1 seed, got %d", s.Quantity(1, 1))
	}
	if n.OfType(1, "seed_earned") == nil {
		t.Error("expected seed_earned notification")
	}
}

func TestTick_MultipleDays_FDAccumulates(t *testing.T) {
	uf := mock.NewUserFlowerRepo()
	w := mock.NewWateringRepo()
	u := mock.NewUserRepo()
	s := mock.NewSeedRepo()
	n := mock.NewNotificationRepo()

	u.Create(context.Background(), "alice", "hash", "Alice")
	uf.Create(context.Background(), 1, 1)

	svc := newTickSvc(uf, w, u, s, n)

	// Simulate 3 days of watering
	base := time.Now().UTC().AddDate(0, 0, -3).Truncate(24 * time.Hour)
	for i := 0; i < 3; i++ {
		d := base.AddDate(0, 0, i)
		w.Record(context.Background(), 1, 1, d)
		svc.RunTick(context.Background(), d)
	}

	user := u.ByID(1)
	// day 1 → +1, day 2 → +2, day 3 → +3 = 6
	if user.FDBalance != 6 {
		t.Errorf("fd_balance: got %d, want 6", user.FDBalance)
	}

	flower := uf.ByID(1)
	if flower.Day != 4 {
		t.Errorf("day: got %d, want 4", flower.Day)
	}
}

func TestTick_SkipsAlreadyDried(t *testing.T) {
	uf := mock.NewUserFlowerRepo()
	w := mock.NewWateringRepo()
	u := mock.NewUserRepo()
	s := mock.NewSeedRepo()
	n := mock.NewNotificationRepo()

	u.Create(context.Background(), "alice", "hash", "Alice")
	uf.Create(context.Background(), 1, 1)
	uf.SetDried(context.Background(), 1)

	svc := newTickSvc(uf, w, u, s, n)
	if err := svc.RunTick(context.Background(), todayDate()); err != nil {
		t.Fatalf("tick error: %v", err)
	}

	user := u.ByID(1)
	if user.FDBalance != 0 {
		t.Error("dried flower should not contribute FD")
	}
}
