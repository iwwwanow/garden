package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/iwwwanow/garden/server/internal/middleware"
	"github.com/iwwwanow/garden/server/internal/repo"
	"github.com/iwwwanow/garden/server/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DevHandler struct {
	tick      service.TickService
	users     repo.UserRepo
	seeds     repo.SeedRepo
	pool      *pgxpool.Pool
	jwtSecret []byte
}

func NewDevHandler(tick service.TickService, users repo.UserRepo, seeds repo.SeedRepo, pool *pgxpool.Pool, jwtSecret string) *DevHandler {
	return &DevHandler{tick: tick, users: users, seeds: seeds, pool: pool, jwtSecret: []byte(jwtSecret)}
}

// Tick simulates midnight: processes today's waterings and advances all active flowers.
// After processing, resets watering state so flowers appear as "needs watering" for the new game day.
func (h *DevHandler) Tick(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	today := time.Now().UTC().Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)

	if err := h.tick.RunTick(ctx, today); err != nil {
		writeError(w, http.StatusInternalServerError, "tick failed")
		return
	}

	// Simulate day rollover: remove today's watering records and roll back last_watered_at
	// so flowers show as "needs watering" for the new game day.
	_, _ = h.pool.Exec(ctx, `DELETE FROM waterings WHERE watered_date = $1`, today)
	_, _ = h.pool.Exec(ctx,
		`UPDATE user_flowers SET last_watered_at = $1 WHERE is_dried = false`, yesterday)

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Users returns all users for switching profiles in dev panel.
func (h *DevHandler) Users(w http.ResponseWriter, r *http.Request) {
	users, err := h.users.ListAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load users")
		return
	}
	writeJSON(w, http.StatusOK, users)
}

// Seeds gives seeds of a flower type to the current user.
// Body: {"flower_id": 1, "quantity": 5}
func (h *DevHandler) Seeds(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FlowerID int `json:"flower_id"`
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.FlowerID == 0 {
		writeError(w, http.StatusBadRequest, "flower_id is required")
		return
	}
	if req.Quantity <= 0 {
		req.Quantity = 5
	}
	userID := middleware.UserIDFromCtx(r.Context())
	if err := h.seeds.Upsert(r.Context(), userID, req.FlowerID, req.Quantity); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to give seeds")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user_id": userID, "flower_id": req.FlowerID, "added": req.Quantity})
}

// Token issues a JWT for any user by ID, no password required. Dev only.
func (h *DevHandler) Token(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID int `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == 0 {
		writeError(w, http.StatusBadRequest, "user_id is required")
		return
	}
	user, err := h.users.GetByID(r.Context(), req.UserID)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	claims := service.Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(h.jwtSecret)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to issue token")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"token": token, "user": user})
}

// Reset truncates all user-generated data (keeps flower templates).
// Use for a clean dev environment without restarting postgres.
func (h *DevHandler) Reset(w http.ResponseWriter, r *http.Request) {
	_, err := h.pool.Exec(r.Context(),
		`TRUNCATE notifications, waterings, seeds, user_flowers, users RESTART IDENTITY CASCADE`,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "reset failed: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "note": "all user data cleared"})
}
