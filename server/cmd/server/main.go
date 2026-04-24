package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iwwwanow/garden/server/config"
	"github.com/iwwwanow/garden/server/internal/handler"
	"github.com/iwwwanow/garden/server/internal/service"
	"github.com/iwwwanow/garden/server/internal/repo"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file, reading from environment")
	}

	cfg := config.Load()

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("ping db: %v", err)
	}
	log.Println("database connected")

	router := handler.NewRouter(cfg, pool)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// daily tick scheduler
	tickSvc := service.NewTickService(
		repo.NewUserFlowerRepo(pool),
		repo.NewWateringRepo(pool),
		repo.NewUserRepo(pool),
		repo.NewSeedRepo(pool),
		repo.NewNotificationRepo(pool),
	)
	go runDailyTick(tickSvc)

	go func() {
		log.Printf("server listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	log.Println("stopped")
}

func runDailyTick(tickSvc service.TickService) {
	for {
		now := time.Now().UTC()
		next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
		time.Sleep(next.Sub(now))

		yesterday := time.Now().UTC().AddDate(0, 0, -1).Truncate(24 * time.Hour)
		if err := tickSvc.RunTick(context.Background(), yesterday); err != nil {
			log.Printf("daily tick error: %v", err)
		}
	}
}
