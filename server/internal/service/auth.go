package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/iwwwanow/garden/server/internal/model"
	"github.com/iwwwanow/garden/server/internal/repo"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthService interface {
	Register(ctx context.Context, username, password, firstName string) (*model.User, error)
	Login(ctx context.Context, username, password string) (token string, user *model.User, err error)
}

type authService struct {
	users     repo.UserRepo
	jwtSecret []byte
}

func NewAuthService(users repo.UserRepo, jwtSecret string) AuthService {
	return &authService{users: users, jwtSecret: []byte(jwtSecret)}
}

func (s *authService) Register(ctx context.Context, username, password, firstName string) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user, err := s.users.Create(ctx, username, string(hash), firstName)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrUsernameTaken
		}
		return nil, err
	}
	return user, nil
}

func (s *authService) Login(ctx context.Context, username, password string) (string, *model.User, error) {
	user, err := s.users.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil, ErrInvalidCredentials
		}
		return "", nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}
	token, err := s.makeToken(user.ID)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (s *authService) makeToken(userID int) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}
