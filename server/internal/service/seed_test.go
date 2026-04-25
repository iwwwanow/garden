package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/iwwwanow/garden/server/internal/repo/mock"
	"github.com/iwwwanow/garden/server/internal/service"
)

func newSeedSvc(s *mock.SeedRepo, u *mock.UserRepo, n *mock.NotificationRepo) service.SeedService {
	return service.NewSeedService(s, u, n)
}

func TestGetSeeds_ReturnsUserSeeds(t *testing.T) {
	s := mock.NewSeedRepo()
	s.Upsert(context.Background(), 1, 1, 3)
	s.Upsert(context.Background(), 1, 2, 5)
	s.Upsert(context.Background(), 2, 1, 1) // другой пользователь

	svc := newSeedSvc(s, mock.NewUserRepo(), mock.NewNotificationRepo())
	seeds, err := svc.GetSeeds(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(seeds) != 2 {
		t.Errorf("expected 2 seed entries for user 1, got %d", len(seeds))
	}
}

func TestShare_Success(t *testing.T) {
	s := mock.NewSeedRepo()
	u := mock.NewUserRepo()
	n := mock.NewNotificationRepo()

	u.Create(context.Background(), "alice", "h", "Alice")
	u.Create(context.Background(), "bob", "h", "Bob")
	s.Upsert(context.Background(), 1, 1, 5)

	svc := newSeedSvc(s, u, n)
	if err := svc.Share(context.Background(), 1, 2, 1, 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if s.Quantity(1, 1) != 2 {
		t.Errorf("sender: expected 2 seeds left, got %d", s.Quantity(1, 1))
	}
	if s.Quantity(2, 1) != 3 {
		t.Errorf("recipient: expected 3 seeds, got %d", s.Quantity(2, 1))
	}
	if n.OfType(2, "seeds_received") == nil {
		t.Error("expected seeds_received notification for recipient")
	}
}

func TestShare_InsufficientSeeds(t *testing.T) {
	s := mock.NewSeedRepo()
	u := mock.NewUserRepo()

	u.Create(context.Background(), "alice", "h", "Alice")
	u.Create(context.Background(), "bob", "h", "Bob")
	s.Upsert(context.Background(), 1, 1, 2)

	svc := newSeedSvc(s, u, mock.NewNotificationRepo())
	err := svc.Share(context.Background(), 1, 2, 1, 5)
	if !errors.Is(err, service.ErrInsufficientSeeds) {
		t.Fatalf("want ErrInsufficientSeeds, got %v", err)
	}
}

func TestShare_NoSeeds(t *testing.T) {
	u := mock.NewUserRepo()
	u.Create(context.Background(), "alice", "h", "Alice")
	u.Create(context.Background(), "bob", "h", "Bob")

	svc := newSeedSvc(mock.NewSeedRepo(), u, mock.NewNotificationRepo())
	err := svc.Share(context.Background(), 1, 2, 1, 1)
	if !errors.Is(err, service.ErrInsufficientSeeds) {
		t.Fatalf("want ErrInsufficientSeeds, got %v", err)
	}
}

func TestShare_RecipientNotFound(t *testing.T) {
	s := mock.NewSeedRepo()
	u := mock.NewUserRepo()

	u.Create(context.Background(), "alice", "h", "Alice")
	s.Upsert(context.Background(), 1, 1, 5)

	svc := newSeedSvc(s, u, mock.NewNotificationRepo())
	err := svc.Share(context.Background(), 1, 999, 1, 1)
	if !errors.Is(err, service.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

func TestShare_InvalidQuantity(t *testing.T) {
	u := mock.NewUserRepo()
	u.Create(context.Background(), "alice", "h", "Alice")
	u.Create(context.Background(), "bob", "h", "Bob")

	svc := newSeedSvc(mock.NewSeedRepo(), u, mock.NewNotificationRepo())
	err := svc.Share(context.Background(), 1, 2, 1, 0)
	if err == nil {
		t.Fatal("expected error for zero quantity")
	}
}
