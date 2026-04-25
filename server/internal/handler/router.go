package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/iwwwanow/garden/server/config"
	mw "github.com/iwwwanow/garden/server/internal/middleware"
	"github.com/iwwwanow/garden/server/internal/repo"
	"github.com/iwwwanow/garden/server/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(cfg *config.Config, pool *pgxpool.Pool) http.Handler {
	// repos
	users := repo.NewUserRepo(pool)
	flowers := repo.NewFlowerRepo(pool)
	userFlowers := repo.NewUserFlowerRepo(pool)
	waterings := repo.NewWateringRepo(pool)
	seeds := repo.NewSeedRepo(pool)
	notifications := repo.NewNotificationRepo(pool)

	// services
	authSvc := service.NewAuthService(users, cfg.JWTSecret)
	flowerSvc := service.NewFlowerService(userFlowers, waterings, seeds, flowers, notifications)
	seedSvc := service.NewSeedService(seeds, users, notifications)
	notifSvc := service.NewNotificationService(notifications)
	tickSvc := service.NewTickService(userFlowers, waterings, users, seeds, notifications)

	// handlers
	authH := NewAuthHandler(authSvc)
	meH := NewMeHandler(users, flowerSvc, seedSvc)
	flowerH := NewFlowerHandler(flowerSvc)
	seedH := NewSeedHandler(seedSvc)
	notifH := NewNotificationHandler(notifSvc)
	lbH := NewLeaderboardHandler(users)
	userH := NewUserHandler(users)
	devH := NewDevHandler(tickSvc, users, seeds, pool)

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// auth (public)
	r.Post("/api/auth/register", authH.Register)
	r.Post("/api/auth/login", authH.Login)

	// protected routes
	r.Group(func(r chi.Router) {
		r.Use(mw.Auth(cfg.JWTSecret))

		r.Get("/api/me", meH.Me)
		r.Put("/api/me", meH.UpdateMe)
		r.Get("/api/leaderboard", lbH.Leaderboard)

		r.Get("/api/flowers", flowerH.List)
		r.Post("/api/flowers/plant", flowerH.Plant)
		r.Post("/api/flowers/{id}/water", flowerH.Water)
		r.Get("/api/flowers/user/{userId}", flowerH.GetUserFlowers)
		r.Get("/api/flowers/{id}", flowerH.GetByID)

		r.Get("/api/seeds", seedH.GetSeeds)
		r.Post("/api/seeds/share", seedH.Share)

		r.Get("/api/notifications", notifH.List)
		r.Patch("/api/notifications/read", notifH.MarkRead)

		r.Get("/api/users/{id}", userH.GetByID)
	})

	// dev endpoints (only in development)
	if cfg.AppEnv == "development" {
		r.Group(func(r chi.Router) {
			r.Use(mw.Auth(cfg.JWTSecret))
			r.Post("/api/dev/tick", devH.Tick)
			r.Get("/api/dev/users", devH.Users)
			r.Post("/api/dev/seeds", devH.Seeds)
			r.Post("/api/dev/reset", devH.Reset)
		})
	}

	return r
}
