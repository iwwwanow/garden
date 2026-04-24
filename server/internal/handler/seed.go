package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/iwwwanow/garden/server/internal/middleware"
	"github.com/iwwwanow/garden/server/internal/service"
)

type SeedHandler struct{ seeds service.SeedService }

func NewSeedHandler(seeds service.SeedService) *SeedHandler { return &SeedHandler{seeds: seeds} }

func (h *SeedHandler) GetSeeds(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromCtx(r.Context())
	seeds, err := h.seeds.GetSeeds(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load seeds")
		return
	}
	writeJSON(w, http.StatusOK, seeds)
}

func (h *SeedHandler) Share(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ToUserID int `json:"to_user_id"`
		FlowerID int `json:"flower_id"`
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ToUserID == 0 || req.FlowerID == 0 || req.Quantity <= 0 {
		writeError(w, http.StatusBadRequest, "to_user_id, flower_id and quantity are required")
		return
	}
	fromUserID := middleware.UserIDFromCtx(r.Context())
	if err := h.seeds.Share(r.Context(), fromUserID, req.ToUserID, req.FlowerID, req.Quantity); err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			writeError(w, http.StatusNotFound, "recipient not found")
		case errors.Is(err, service.ErrInsufficientSeeds):
			writeError(w, http.StatusUnprocessableEntity, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "share failed")
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
