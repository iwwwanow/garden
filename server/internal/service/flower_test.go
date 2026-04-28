package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/iwwwanow/garden/server/internal/model"
	"github.com/iwwwanow/garden/server/internal/repo/mock"
	"github.com/iwwwanow/garden/server/internal/service"
)

func newFlowerSvc(
	uf *mock.UserFlowerRepo,
	w *mock.WateringRepo,
	s *mock.SeedRepo,
	f *mock.FlowerRepo,
	n *mock.NotificationRepo,
) service.FlowerService {
	return service.NewFlowerService(uf, w, s, f, n)
}

// ── Water ────────────────────────────────────────────────────────────────────

func TestWater_Success(t *testing.T) {
	uf := mock.NewUserFlowerRepo()
	w := mock.NewWateringRepo()
	flowers := mock.NewFlowerRepo(&model.Flower{ID: 1, Season: 1, ImagePath: "f.png"})
	n := mock.NewNotificationRepo()

	uf.Create(context.Background(), 1, 1) // user 1 owns flower #1

	svc := newFlowerSvc(uf, w, mock.NewSeedRepo(), flowers, n)
	if err := svc.Water(context.Background(), 1, 2); err != nil { // user 2 waters
		t.Fatalf("unexpected error: %v", err)
	}

	// notification sent to owner
	if n.OfType(1, "flower_watered") == nil {
		t.Error("expected flower_watered notification for owner")
	}
}

func TestWater_SelfNoNotification(t *testing.T) {
	uf := mock.NewUserFlowerRepo()
	w := mock.NewWateringRepo()
	n := mock.NewNotificationRepo()

	uf.Create(context.Background(), 1, 1)

	svc := newFlowerSvc(uf, w, mock.NewSeedRepo(), mock.NewFlowerRepo(), n)
	svc.Water(context.Background(), 1, 1) // owner waters own flower

	if n.OfType(1, "flower_watered") != nil {
		t.Error("no notification expected when owner waters own flower")
	}
}

func TestWater_AlreadyWatered(t *testing.T) {
	uf := mock.NewUserFlowerRepo()
	w := mock.NewWateringRepo()
	n := mock.NewNotificationRepo()

	uf.Create(context.Background(), 1, 1)
	svc := newFlowerSvc(uf, w, mock.NewSeedRepo(), mock.NewFlowerRepo(), n)

	svc.Water(context.Background(), 1, 2) // first time

	err := svc.Water(context.Background(), 1, 2) // second time same day
	if !errors.Is(err, service.ErrAlreadyWatered) {
		t.Fatalf("want ErrAlreadyWatered, got %v", err)
	}
}

func TestWater_DriedFlower(t *testing.T) {
	uf := mock.NewUserFlowerRepo()
	w := mock.NewWateringRepo()
	n := mock.NewNotificationRepo()

	uf.Create(context.Background(), 1, 1)
	uf.SetDried(context.Background(), 1)

	svc := newFlowerSvc(uf, w, mock.NewSeedRepo(), mock.NewFlowerRepo(), n)
	err := svc.Water(context.Background(), 1, 2)
	if !errors.Is(err, service.ErrFlowerDried) {
		t.Fatalf("want ErrFlowerDried, got %v", err)
	}
}

func TestWater_FlowerNotFound(t *testing.T) {
	svc := newFlowerSvc(
		mock.NewUserFlowerRepo(), mock.NewWateringRepo(),
		mock.NewSeedRepo(), mock.NewFlowerRepo(), mock.NewNotificationRepo(),
	)
	err := svc.Water(context.Background(), 999, 1)
	if !errors.Is(err, service.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

// ── Plant ────────────────────────────────────────────────────────────────────

func TestPlant_Success(t *testing.T) {
	uf := mock.NewUserFlowerRepo()
	s := mock.NewSeedRepo()
	flowers := mock.NewFlowerRepo(&model.Flower{ID: 1, Season: 1, ImagePath: "f.png"})
	n := mock.NewNotificationRepo()

	s.Upsert(context.Background(), 1, 1, 3) // give user 1 some seeds

	svc := newFlowerSvc(uf, mock.NewWateringRepo(), s, flowers, n)
	planted, err := svc.Plant(context.Background(), 1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if planted.UserID != 1 || planted.FlowerID != 1 {
		t.Errorf("wrong planted flower: %+v", planted)
	}
	if s.Quantity(1, 1) != 2 {
		t.Errorf("expected 2 seeds left, got %d", s.Quantity(1, 1))
	}
}

func TestPlant_InsufficientSeeds(t *testing.T) {
	flowers := mock.NewFlowerRepo(&model.Flower{ID: 1})
	svc := newFlowerSvc(
		mock.NewUserFlowerRepo(), mock.NewWateringRepo(),
		mock.NewSeedRepo(), flowers, mock.NewNotificationRepo(),
	)
	_, err := svc.Plant(context.Background(), 1, 1)
	if !errors.Is(err, service.ErrInsufficientSeeds) {
		t.Fatalf("want ErrInsufficientSeeds, got %v", err)
	}
}

func TestPlant_MaxFlowers(t *testing.T) {
	uf := mock.NewUserFlowerRepo()
	s := mock.NewSeedRepo()
	flowers := mock.NewFlowerRepo(&model.Flower{ID: 1})
	n := mock.NewNotificationRepo()

	// fill up to limit
	for i := 0; i < 64; i++ {
		uf.Create(context.Background(), 1, 1)
	}
	s.Upsert(context.Background(), 1, 1, 10)

	svc := newFlowerSvc(uf, mock.NewWateringRepo(), s, flowers, n)
	_, err := svc.Plant(context.Background(), 1, 1)
	if !errors.Is(err, service.ErrMaxFlowers) {
		t.Fatalf("want ErrMaxFlowers, got %v", err)
	}
}

func TestPlant_DailyLimitReached(t *testing.T) {
	uf := mock.NewUserFlowerRepo()
	s := mock.NewSeedRepo()
	flowers := mock.NewFlowerRepo(&model.Flower{ID: 1})
	n := mock.NewNotificationRepo()

	s.Upsert(context.Background(), 1, 1, 5)
	uf.Create(context.Background(), 1, 1) // already planted one today

	svc := newFlowerSvc(uf, mock.NewWateringRepo(), s, flowers, n)
	_, err := svc.Plant(context.Background(), 1, 1)
	if !errors.Is(err, service.ErrPlantLimitReached) {
		t.Fatalf("want ErrPlantLimitReached, got %v", err)
	}
}

func TestPlant_FlowerTemplateNotFound(t *testing.T) {
	s := mock.NewSeedRepo()
	s.Upsert(context.Background(), 1, 99, 1)

	svc := newFlowerSvc(
		mock.NewUserFlowerRepo(), mock.NewWateringRepo(),
		s, mock.NewFlowerRepo(), mock.NewNotificationRepo(),
	)
	_, err := svc.Plant(context.Background(), 1, 99)
	if !errors.Is(err, service.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

// ── ListFlowers ───────────────────────────────────────────────────────────────

func TestListFlowers_ReturnsAll(t *testing.T) {
	flowers := mock.NewFlowerRepo(
		&model.Flower{ID: 1, Season: 1, ImagePath: "a.png"},
		&model.Flower{ID: 2, Season: 1, ImagePath: "b.png"},
	)
	svc := newFlowerSvc(mock.NewUserFlowerRepo(), mock.NewWateringRepo(), mock.NewSeedRepo(), flowers, mock.NewNotificationRepo())

	list, err := svc.ListFlowers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 flowers, got %d", len(list))
	}
}

func TestListFlowers_Empty(t *testing.T) {
	svc := newFlowerSvc(mock.NewUserFlowerRepo(), mock.NewWateringRepo(), mock.NewSeedRepo(), mock.NewFlowerRepo(), mock.NewNotificationRepo())

	list, err := svc.ListFlowers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d", len(list))
	}
}

// ── GetFlower ─────────────────────────────────────────────────────────────────

func TestGetFlower_Found(t *testing.T) {
	flowers := mock.NewFlowerRepo(&model.Flower{ID: 1, Season: 1, ImagePath: "a.png"})
	svc := newFlowerSvc(mock.NewUserFlowerRepo(), mock.NewWateringRepo(), mock.NewSeedRepo(), flowers, mock.NewNotificationRepo())

	f, err := svc.GetFlower(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.ID != 1 {
		t.Errorf("expected flower ID 1, got %d", f.ID)
	}
}

func TestGetFlower_NotFound(t *testing.T) {
	svc := newFlowerSvc(mock.NewUserFlowerRepo(), mock.NewWateringRepo(), mock.NewSeedRepo(), mock.NewFlowerRepo(), mock.NewNotificationRepo())

	_, err := svc.GetFlower(context.Background(), 99)
	if !errors.Is(err, service.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func todayDate() time.Time {
	return time.Now().UTC().Truncate(24 * time.Hour)
}
