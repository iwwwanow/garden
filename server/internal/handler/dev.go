package handler

import (
	"net/http"
	"time"

	"github.com/iwwwanow/garden/server/internal/repo"
	"github.com/iwwwanow/garden/server/internal/service"
)

type DevHandler struct {
	tick  service.TickService
	users repo.UserRepo
}

func NewDevHandler(tick service.TickService, users repo.UserRepo) *DevHandler {
	return &DevHandler{tick: tick, users: users}
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
