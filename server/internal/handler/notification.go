package handler

import (
	"net/http"

	"github.com/iwwwanow/garden/server/internal/middleware"
	"github.com/iwwwanow/garden/server/internal/service"
)

type NotificationHandler struct{ notifications service.NotificationService }

func NewNotificationHandler(n service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifications: n}
}

func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromCtx(r.Context())
	list, err := h.notifications.List(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load notifications")
		return
	}
	writeJSON(w, http.StatusOK, list)
}
