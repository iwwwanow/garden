package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/iwwwanow/garden/server/internal/middleware"
	"github.com/iwwwanow/garden/server/internal/repo"
	"github.com/iwwwanow/garden/server/internal/service"
)

type MeHandler struct {
	users   repo.UserRepo
	flowers service.FlowerService
	seeds   service.SeedService
}

func NewMeHandler(users repo.UserRepo, flowers service.FlowerService, seeds service.SeedService) *MeHandler {
	return &MeHandler{users: users, flowers: flowers, seeds: seeds}
}

func (h *MeHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FirstName string `json:"first_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.FirstName = strings.TrimSpace(req.FirstName)
	if req.FirstName == "" {
		writeError(w, http.StatusBadRequest, "first_name is required")
		return
	}
	userID := middleware.UserIDFromCtx(r.Context())
	if err := h.users.UpdateFirstName(r.Context(), userID, req.FirstName); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update profile")
		return
	}
	user, err := h.users.GetByID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load user")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *MeHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromCtx(r.Context())
	ctx := r.Context()

	user, err := h.users.GetByID(ctx, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load user")
		return
	}
	flowers, err := h.flowers.GetUserFlowers(ctx, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load flowers")
		return
	}
	seeds, err := h.seeds.GetSeeds(ctx, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load seeds")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"user":    user,
		"flowers": flowers,
		"seeds":   seeds,
	})
}
