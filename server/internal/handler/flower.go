package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iwwwanow/garden/server/internal/middleware"
	"github.com/iwwwanow/garden/server/internal/service"
)

type FlowerHandler struct{ flowers service.FlowerService }

func NewFlowerHandler(flowers service.FlowerService) *FlowerHandler {
	return &FlowerHandler{flowers: flowers}
}

func (h *FlowerHandler) List(w http.ResponseWriter, r *http.Request) {
	flowers, err := h.flowers.ListFlowers(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load flowers")
		return
	}
	writeJSON(w, http.StatusOK, flowers)
}

func (h *FlowerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid flower id")
		return
	}
	flower, err := h.flowers.GetFlower(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "flower not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to load flower")
		return
	}
	writeJSON(w, http.StatusOK, flower)
}

func (h *FlowerHandler) GetUserFlowers(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	flowers, err := h.flowers.GetUserFlowers(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load flowers")
		return
	}
	writeJSON(w, http.StatusOK, flowers)
}

func (h *FlowerHandler) Water(w http.ResponseWriter, r *http.Request) {
	userFlowerID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid flower id")
		return
	}
	callerID := middleware.UserIDFromCtx(r.Context())
	if err := h.flowers.Water(r.Context(), userFlowerID, callerID); err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, service.ErrFlowerDried):
			writeError(w, http.StatusUnprocessableEntity, err.Error())
		case errors.Is(err, service.ErrAlreadyWatered):
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "watering failed")
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *FlowerHandler) Plant(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FlowerID int `json:"flower_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.FlowerID == 0 {
		writeError(w, http.StatusBadRequest, "flower_id is required")
		return
	}
	userID := middleware.UserIDFromCtx(r.Context())
	uf, err := h.flowers.Plant(r.Context(), userID, req.FlowerID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			writeError(w, http.StatusNotFound, "flower template not found")
		case errors.Is(err, service.ErrMaxFlowers):
			writeError(w, http.StatusUnprocessableEntity, err.Error())
		case errors.Is(err, service.ErrInsufficientSeeds):
			writeError(w, http.StatusUnprocessableEntity, err.Error())
		case errors.Is(err, service.ErrPlantLimitReached):
			writeError(w, http.StatusUnprocessableEntity, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "planting failed")
		}
		return
	}
	writeJSON(w, http.StatusCreated, uf)
}
