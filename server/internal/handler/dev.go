package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/iwwwanow/garden/server/internal/middleware"
	"github.com/iwwwanow/garden/server/internal/repo"
	"github.com/iwwwanow/garden/server/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DevHandler struct {
	tick  service.TickService
	users repo.UserRepo
	seeds repo.SeedRepo
	pool  *pgxpool.Pool
}

func NewDevHandler(tick service.TickService, users repo.UserRepo, seeds repo.SeedRepo, pool *pgxpool.Pool) *DevHandler {
	return &DevHandler{tick: tick, users: users, seeds: seeds, pool: pool}
}

// Tick simulates midnight: processes today's waterings and advances all active flowers.
func (h *DevHandler) Tick(w http.ResponseWriter, r *http.Request) {
	today := time.Now().UTC().Truncate(24 * time.Hour)
	if err := h.tick.RunTick(r.Context(), today); err != nil {
		writeError(w, http.StatusInternalServerError, "tick failed")
		return
	}
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
