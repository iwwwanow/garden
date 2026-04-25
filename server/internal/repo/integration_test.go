package repo_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/iwwwanow/garden/server/internal/repo"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // register "pgx" driver for database/sql
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// migrationsDir returns the absolute path to the migrations directory.
func migrationsDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "migrations")
}

// setupDB starts a postgres container, runs migrations, and returns a pool.
// The container is terminated when the test finishes.
func setupDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	ctr, err := postgres.Run(ctx,
		"postgres:17.4-alpine",
		postgres.WithDatabase("gardentest"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}
	t.Cleanup(func() { ctr.Terminate(ctx) }) //nolint:errcheck

	connStr, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("connection string: %v", err)
	}

	// run migrations via goose
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatalf("open sql.DB: %v", err)
	}
	defer db.Close()

	goose.SetLogger(goose.NopLogger())
	if err := goose.Up(db, migrationsDir()); err != nil {
		t.Fatalf("goose up: %v", err)
	}

	// connect via pgxpool
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("pgxpool new: %v", err)
	}
	t.Cleanup(pool.Close)

	return pool
}

// ── UserRepo ─────────────────────────────────────────────────────────────────

func TestUserRepo_CreateAndGet(t *testing.T) {
	pool := setupDB(t)
	ctx := context.Background()
	r := repo.NewUserRepo(pool)

	user, err := r.Create(ctx, "alice", "hash", "Alice")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if user.ID == 0 {
		t.Error("expected non-zero ID")
	}

	got, err := r.GetByUsername(ctx, "alice")
	if err != nil {
		t.Fatalf("GetByUsername: %v", err)
	}
	if got.ID != user.ID {
		t.Errorf("ID mismatch: got %d, want %d", got.ID, user.ID)
	}

	byID, err := r.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if byID.Username != "alice" {
		t.Errorf("username: got %q", byID.Username)
	}
}

func TestUserRepo_AddFDAndLeaderboard(t *testing.T) {
	pool := setupDB(t)
	ctx := context.Background()
	r := repo.NewUserRepo(pool)

	alice, _ := r.Create(ctx, "alice", "hash", "Alice")
	bob, _ := r.Create(ctx, "bob", "hash", "Bob")

	r.AddFD(ctx, alice.ID, 100)
	r.AddFD(ctx, bob.ID, 50)

	lb, err := r.Leaderboard(ctx, 10, 0)
	if err != nil {
		t.Fatalf("Leaderboard: %v", err)
	}
	if len(lb) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(lb))
	}
	if lb[0].Username != "alice" {
		t.Errorf("top user: got %q, want alice", lb[0].Username)
	}
}

// ── Full flow: register → plant → water → tick ────────────────────────────────

func TestFullFlow(t *testing.T) {
	pool := setupDB(t)
	ctx := context.Background()

	users := repo.NewUserRepo(pool)
	userFlowers := repo.NewUserFlowerRepo(pool)
	waterings := repo.NewWateringRepo(pool)
	seeds := repo.NewSeedRepo(pool)
	notifications := repo.NewNotificationRepo(pool)

	// seed a flower template
	_, err := pool.Exec(ctx,
		`INSERT INTO flowers (season, image_path) VALUES (1, 'test.png')`,
	)
	if err != nil {
		t.Fatalf("insert flower: %v", err)
	}

	// register user
	user, err := users.Create(ctx, "testuser", "hash", "Test")
	if err != nil {
		t.Fatalf("Create user: %v", err)
	}

	// give a seed and plant
	seeds.Upsert(ctx, user.ID, 1, 3)
	uf, err := userFlowers.Create(ctx, user.ID, 1)
	if err != nil {
		t.Fatalf("Create user_flower: %v", err)
	}

	// water the flower
	today := time.Now().UTC().Truncate(24 * time.Hour)
	if err := waterings.Record(ctx, uf.ID, user.ID, today); err != nil {
		t.Fatalf("Record watering: %v", err)
	}

	// verify watering recorded
	watered, err := waterings.AnyOnDate(ctx, uf.ID, today)
	if err != nil {
		t.Fatalf("AnyOnDate: %v", err)
	}
	if !watered {
		t.Error("expected watering to be recorded")
	}

	// simulate tick: award FD, increment day
	users.AddFD(ctx, user.ID, uf.Day)
	userFlowers.IncrementDay(ctx, uf.ID)

	// verify results
	updatedUser, _ := users.GetByID(ctx, user.ID)
	if updatedUser.FDBalance != 1 {
		t.Errorf("fd_balance: got %d, want 1", updatedUser.FDBalance)
	}

	updatedFlower, _ := userFlowers.GetByID(ctx, uf.ID)
	if updatedFlower.Day != 2 {
		t.Errorf("day: got %d, want 2", updatedFlower.Day)
	}

	// create notification
	notifications.Create(ctx, user.ID, "test_event", map[string]any{"foo": "bar"})
	notifs, err := notifications.ListByUserID(ctx, user.ID)
	if err != nil {
		t.Fatalf("ListByUserID: %v", err)
	}
	if len(notifs) != 1 {
		t.Errorf("expected 1 notification, got %d", len(notifs))
	}
}
