// Package mock provides in-memory implementations of all repo interfaces for unit testing.
package mock

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/iwwwanow/garden/server/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ── UserRepo ─────────────────────────────────────────────────────────────────

type UserRepo struct {
	users  []*model.User
	nextID int
}

func NewUserRepo() *UserRepo { return &UserRepo{nextID: 1} }

func (r *UserRepo) Create(_ context.Context, username, passwordHash, firstName string) (*model.User, error) {
	for _, u := range r.users {
		if u.Username == username {
			return nil, &pgconn.PgError{Code: "23505"}
		}
	}
	u := &model.User{
		ID: r.nextID, Username: username, PasswordHash: passwordHash,
		FirstName: firstName, CreatedAt: time.Now(),
	}
	r.nextID++
	r.users = append(r.users, u)
	return copyUser(u), nil
}

func (r *UserRepo) GetByUsername(_ context.Context, username string) (*model.User, error) {
	for _, u := range r.users {
		if u.Username == username {
			return copyUser(u), nil
		}
	}
	return nil, pgx.ErrNoRows
}

func (r *UserRepo) GetByID(_ context.Context, id int) (*model.User, error) {
	for _, u := range r.users {
		if u.ID == id {
			return copyUser(u), nil
		}
	}
	return nil, pgx.ErrNoRows
}

func (r *UserRepo) AddFD(_ context.Context, userID, delta int) error {
	for _, u := range r.users {
		if u.ID == userID {
			u.FDBalance += delta
			return nil
		}
	}
	return pgx.ErrNoRows
}

func (r *UserRepo) UpdateFirstName(_ context.Context, userID int, firstName string) error {
	for _, u := range r.users {
		if u.ID == userID {
			u.FirstName = firstName
			return nil
		}
	}
	return pgx.ErrNoRows
}

func (r *UserRepo) Leaderboard(_ context.Context, _, _ int) ([]*model.User, error) {
	result := make([]*model.User, len(r.users))
	for i, u := range r.users {
		result[i] = copyUser(u)
	}
	return result, nil
}

func (r *UserRepo) ListAll(_ context.Context) ([]*model.User, error) {
	return r.Leaderboard(context.Background(), len(r.users), 0)
}

// ByID is a test-only helper to inspect stored state.
func (r *UserRepo) ByID(id int) *model.User {
	for _, u := range r.users {
		if u.ID == id {
			return u
		}
	}
	return nil
}

func copyUser(u *model.User) *model.User { c := *u; return &c }

// ── FlowerRepo ────────────────────────────────────────────────────────────────

type FlowerRepo struct {
	flowers []*model.Flower
}

func NewFlowerRepo(flowers ...*model.Flower) *FlowerRepo { return &FlowerRepo{flowers: flowers} }

func (r *FlowerRepo) GetByID(_ context.Context, id int) (*model.Flower, error) {
	for _, f := range r.flowers {
		if f.ID == id {
			return f, nil
		}
	}
	return nil, pgx.ErrNoRows
}

func (r *FlowerRepo) ListAll(_ context.Context) ([]*model.Flower, error) { return r.flowers, nil }

// ── UserFlowerRepo ────────────────────────────────────────────────────────────

type UserFlowerRepo struct {
	flowers []*model.UserFlower
	nextID  int
}

func NewUserFlowerRepo() *UserFlowerRepo { return &UserFlowerRepo{nextID: 1} }

func (r *UserFlowerRepo) Create(_ context.Context, userID, flowerID int) (*model.UserFlower, error) {
	uf := &model.UserFlower{
		ID: r.nextID, UserID: userID, FlowerID: flowerID, Day: 1, CreatedAt: time.Now(),
	}
	r.nextID++
	r.flowers = append(r.flowers, uf)
	return uf, nil
}

func (r *UserFlowerRepo) GetByID(_ context.Context, id int) (*model.UserFlower, error) {
	for _, uf := range r.flowers {
		if uf.ID == id {
			c := *uf
			return &c, nil
		}
	}
	return nil, pgx.ErrNoRows
}

func (r *UserFlowerRepo) ListByUserID(_ context.Context, userID int) ([]*model.UserFlower, error) {
	var result []*model.UserFlower
	for _, uf := range r.flowers {
		if uf.UserID == userID {
			c := *uf
			result = append(result, &c)
		}
	}
	return result, nil
}

func (r *UserFlowerRepo) CountActiveByUserID(_ context.Context, userID int) (int, error) {
	count := 0
	for _, uf := range r.flowers {
		if uf.UserID == userID && !uf.IsDried {
			count++
		}
	}
	return count, nil
}

func (r *UserFlowerRepo) ListActive(_ context.Context) ([]*model.UserFlower, error) {
	var result []*model.UserFlower
	for _, uf := range r.flowers {
		if !uf.IsDried {
			c := *uf
			result = append(result, &c)
		}
	}
	return result, nil
}

func (r *UserFlowerRepo) IncrementDay(_ context.Context, id int) error {
	for _, uf := range r.flowers {
		if uf.ID == id {
			uf.Day++
			return nil
		}
	}
	return pgx.ErrNoRows
}

func (r *UserFlowerRepo) SetDried(_ context.Context, id int) error {
	for _, uf := range r.flowers {
		if uf.ID == id {
			uf.IsDried = true
			return nil
		}
	}
	return pgx.ErrNoRows
}

func (r *UserFlowerRepo) SetLastWatered(_ context.Context, id int) error {
	now := time.Now()
	for _, uf := range r.flowers {
		if uf.ID == id {
			uf.LastWateredAt = &now
			return nil
		}
	}
	return pgx.ErrNoRows
}

// ByID is a test-only helper.
func (r *UserFlowerRepo) ByID(id int) *model.UserFlower {
	for _, uf := range r.flowers {
		if uf.ID == id {
			return uf
		}
	}
	return nil
}

// ── WateringRepo ──────────────────────────────────────────────────────────────

type WateringRepo struct {
	records []wateringRecord
}

type wateringRecord struct {
	userFlowerID    int
	wateredByUserID int
	date            time.Time
}

func NewWateringRepo() *WateringRepo { return &WateringRepo{} }

func (r *WateringRepo) Record(_ context.Context, userFlowerID, wateredByUserID int, date time.Time) error {
	// unique per (user_flower_id, date) — one watering per flower per day
	for _, w := range r.records {
		if w.userFlowerID == userFlowerID && w.date.Equal(date) {
			return &pgconn.PgError{Code: "23505"}
		}
	}
	r.records = append(r.records, wateringRecord{userFlowerID, wateredByUserID, date})
	return nil
}

func (r *WateringRepo) HasWatered(_ context.Context, userFlowerID, wateredByUserID int, date time.Time) (bool, error) {
	for _, w := range r.records {
		if w.userFlowerID == userFlowerID && w.wateredByUserID == wateredByUserID && w.date.Equal(date) {
			return true, nil
		}
	}
	return false, nil
}

func (r *WateringRepo) AnyOnDate(_ context.Context, userFlowerID int, date time.Time) (bool, error) {
	for _, w := range r.records {
		if w.userFlowerID == userFlowerID && w.date.Equal(date) {
			return true, nil
		}
	}
	return false, nil
}

// ── SeedRepo ──────────────────────────────────────────────────────────────────

type SeedRepo struct {
	seeds  map[string]*model.Seed
	nextID int
}

func NewSeedRepo() *SeedRepo { return &SeedRepo{seeds: make(map[string]*model.Seed), nextID: 1} }

func seedKey(userID, flowerID int) string { return fmt.Sprintf("%d:%d", userID, flowerID) }

func (r *SeedRepo) ListByUserID(_ context.Context, userID int) ([]*model.Seed, error) {
	var result []*model.Seed
	for _, s := range r.seeds {
		if s.UserID == userID {
			c := *s
			result = append(result, &c)
		}
	}
	return result, nil
}

func (r *SeedRepo) GetByUserAndFlower(_ context.Context, userID, flowerID int) (*model.Seed, error) {
	if s, ok := r.seeds[seedKey(userID, flowerID)]; ok {
		c := *s
		return &c, nil
	}
	return nil, pgx.ErrNoRows
}

func (r *SeedRepo) Upsert(_ context.Context, userID, flowerID, delta int) error {
	k := seedKey(userID, flowerID)
	if s, ok := r.seeds[k]; ok {
		s.Quantity += delta
	} else {
		r.seeds[k] = &model.Seed{ID: r.nextID, UserID: userID, FlowerID: flowerID, Quantity: delta}
		r.nextID++
	}
	return nil
}

// Quantity is a test-only helper.
func (r *SeedRepo) Quantity(userID, flowerID int) int {
	if s, ok := r.seeds[seedKey(userID, flowerID)]; ok {
		return s.Quantity
	}
	return 0
}

// ── NotificationRepo ──────────────────────────────────────────────────────────

type NotificationRepo struct {
	Notifications []*model.Notification
	nextID        int
}

func NewNotificationRepo() *NotificationRepo { return &NotificationRepo{nextID: 1} }

func (r *NotificationRepo) MarkAllRead(_ context.Context, userID int) error {
	for _, n := range r.Notifications {
		if n.UserID == userID {
			n.IsRead = true
		}
	}
	return nil
}

func (r *NotificationRepo) Create(_ context.Context, userID int, notifType string, payload map[string]any) error {
	data, _ := json.Marshal(payload)
	r.Notifications = append(r.Notifications, &model.Notification{
		ID: r.nextID, UserID: userID, Type: notifType, Payload: data, CreatedAt: time.Now(),
	})
	r.nextID++
	return nil
}

func (r *NotificationRepo) ListByUserID(_ context.Context, userID int) ([]*model.Notification, error) {
	var result []*model.Notification
	for _, n := range r.Notifications {
		if n.UserID == userID {
			c := *n
			result = append(result, &c)
		}
	}
	return result, nil
}

// OfType returns the first notification of the given type for userID.
func (r *NotificationRepo) OfType(userID int, notifType string) *model.Notification {
	for _, n := range r.Notifications {
		if n.UserID == userID && n.Type == notifType {
			return n
		}
	}
	return nil
}
