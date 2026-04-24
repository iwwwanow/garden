package handler

import (
	"net/http"

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
