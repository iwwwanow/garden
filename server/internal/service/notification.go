package service

import (
	"context"

	"github.com/iwwwanow/garden/server/internal/model"
	"github.com/iwwwanow/garden/server/internal/repo"
)

type NotificationService interface {
	List(ctx context.Context, userID int) ([]*model.Notification, error)
}

type notificationService struct {
	notifications repo.NotificationRepo
}

func NewNotificationService(n repo.NotificationRepo) NotificationService {
	return &notificationService{notifications: n}
}

func (s *notificationService) List(ctx context.Context, userID int) ([]*model.Notification, error) {
	return s.notifications.ListByUserID(ctx, userID)
}
