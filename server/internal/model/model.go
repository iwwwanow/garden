package model

import (
	"encoding/json"
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	FirstName    string    `json:"first_name"`
	TelegramID   *string   `json:"telegram_id,omitempty"`
	FDBalance    int       `json:"fd_balance"`
	CreatedAt    time.Time `json:"created_at"`
}

type Flower struct {
	ID        int       `json:"id"`
	Season    int       `json:"season"`
	ImagePath string    `json:"image_path"`
	CreatedAt time.Time `json:"created_at"`
}

type UserFlower struct {
	ID            int        `json:"id"`
	UserID        int        `json:"user_id"`
	FlowerID      int        `json:"flower_id"`
	Day           int        `json:"day"`
	LastWateredAt *time.Time `json:"last_watered_at,omitempty"`
	IsDried       bool       `json:"is_dried"`
	CreatedAt     time.Time  `json:"created_at"`
}

type Seed struct {
	ID       int `json:"id"`
	UserID   int `json:"user_id"`
	FlowerID int `json:"flower_id"`
	Quantity int `json:"quantity"`
}

type Notification struct {
	ID        int             `json:"id"`
	UserID    int             `json:"user_id"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	IsRead    bool            `json:"is_read"`
	CreatedAt time.Time       `json:"created_at"`
}
