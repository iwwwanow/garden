package handler

import (
	"net/http"
	"strconv"

	"github.com/iwwwanow/garden/server/internal/repo"
)

type LeaderboardHandler struct{ users repo.UserRepo }

func NewLeaderboardHandler(users repo.UserRepo) *LeaderboardHandler {
	return &LeaderboardHandler{users: users}
}

func (h *LeaderboardHandler) Leaderboard(w http.ResponseWriter, r *http.Request) {
	limit := 100
	offset := 0
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}
	users, err := h.users.Leaderboard(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load leaderboard")
		return
	}
	writeJSON(w, http.StatusOK, users)
}
