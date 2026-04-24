package handler

import (
	"net/http"

	"github.com/iwwwanow/garden/server/internal/repo"
)

type LeaderboardHandler struct{ users repo.UserRepo }

func NewLeaderboardHandler(users repo.UserRepo) *LeaderboardHandler {
	return &LeaderboardHandler{users: users}
}

func (h *LeaderboardHandler) Leaderboard(w http.ResponseWriter, r *http.Request) {
	users, err := h.users.Leaderboard(r.Context(), 100)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load leaderboard")
		return
	}
	writeJSON(w, http.StatusOK, users)
}
