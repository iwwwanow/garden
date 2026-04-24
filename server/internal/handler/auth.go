package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/iwwwanow/garden/server/internal/service"
)

type AuthHandler struct{ auth service.AuthService }

func NewAuthHandler(auth service.AuthService) *AuthHandler { return &AuthHandler{auth: auth} }

type registerRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Username == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "username and password are required")
		return
	}
	user, err := h.auth.Register(r.Context(), req.Username, req.Password, req.FirstName)
	if err != nil {
		if errors.Is(err, service.ErrUsernameTaken) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "registration failed")
		return
	}
	writeJSON(w, http.StatusCreated, user)
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	token, user, err := h.auth.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "login failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"token": token, "user": user})
}
