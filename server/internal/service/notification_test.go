package service_test

import (
	"context"
	"testing"

	"github.com/iwwwanow/garden/server/internal/repo/mock"
	"github.com/iwwwanow/garden/server/internal/service"
)

func newNotifSvc(n *mock.NotificationRepo) service.NotificationService {
	return service.NewNotificationService(n)
}

func TestNotification_List(t *testing.T) {
	n := mock.NewNotificationRepo()
	n.Create(context.Background(), 1, "flower_watered", map[string]any{"x": 1})
	n.Create(context.Background(), 1, "seed_earned", map[string]any{"x": 2})
	n.Create(context.Background(), 2, "flower_watered", map[string]any{"x": 3}) // другой пользователь

	svc := newNotifSvc(n)
	list, err := svc.List(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 notifications for user 1, got %d", len(list))
	}
}

func TestNotification_List_Empty(t *testing.T) {
	svc := newNotifSvc(mock.NewNotificationRepo())
	list, err := svc.List(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d", len(list))
	}
}

func TestNotification_MarkAllRead(t *testing.T) {
	n := mock.NewNotificationRepo()
	n.Create(context.Background(), 1, "flower_watered", map[string]any{})
	n.Create(context.Background(), 1, "seed_earned", map[string]any{})
	n.Create(context.Background(), 2, "flower_watered", map[string]any{}) // другой пользователь

	svc := newNotifSvc(n)
	if err := svc.MarkAllRead(context.Background(), 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list, _ := svc.List(context.Background(), 1)
	for _, notif := range list {
		if !notif.IsRead {
			t.Errorf("notification %d should be read", notif.ID)
		}
	}

	// уведомления другого пользователя не затронуты
	list2, _ := svc.List(context.Background(), 2)
	for _, notif := range list2 {
		if notif.IsRead {
			t.Errorf("notification %d of user 2 should not be read", notif.ID)
		}
	}
}
