package service_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/iwwwanow/garden/server/internal/repo/mock"
	"github.com/iwwwanow/garden/server/internal/service"
)

func TestRegister_Success(t *testing.T) {
	svc := service.NewAuthService(mock.NewUserRepo(), "secret")

	user, err := svc.Register(context.Background(), "alice", "password", "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Username != "alice" {
		t.Errorf("username: got %q, want %q", user.Username, "alice")
	}
	if user.ID == 0 {
		t.Error("expected non-zero ID")
	}
	// PasswordHash must not appear in JSON output (json:"-")
	data, _ := json.Marshal(user)
	if bytes.Contains(data, []byte("password_hash")) {
		t.Error("password_hash must not appear in JSON response")
	}
}

func TestRegister_DuplicateUsername(t *testing.T) {
	users := mock.NewUserRepo()
	svc := service.NewAuthService(users, "secret")

	if _, err := svc.Register(context.Background(), "alice", "pass", "Alice"); err != nil {
		t.Fatalf("first register: %v", err)
	}
	_, err := svc.Register(context.Background(), "alice", "other", "Other")
	if !errors.Is(err, service.ErrUsernameTaken) {
		t.Fatalf("want ErrUsernameTaken, got %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	users := mock.NewUserRepo()
	svc := service.NewAuthService(users, "secret")

	if _, err := svc.Register(context.Background(), "alice", "pass123", "Alice"); err != nil {
		t.Fatalf("register: %v", err)
	}

	token, user, err := svc.Login(context.Background(), "alice", "pass123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
	if user.Username != "alice" {
		t.Errorf("username: got %q", user.Username)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	users := mock.NewUserRepo()
	svc := service.NewAuthService(users, "secret")

	svc.Register(context.Background(), "alice", "pass123", "Alice") //nolint

	_, _, err := svc.Login(context.Background(), "alice", "wrongpass")
	if !errors.Is(err, service.ErrInvalidCredentials) {
		t.Fatalf("want ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	svc := service.NewAuthService(mock.NewUserRepo(), "secret")

	_, _, err := svc.Login(context.Background(), "nobody", "pass")
	if !errors.Is(err, service.ErrInvalidCredentials) {
		t.Fatalf("want ErrInvalidCredentials, got %v", err)
	}
}
