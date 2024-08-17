package testhelpers

import (
	"context"
	"log"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	migrationsDir    = "../migrations"
	migrationsSuffix = "up.sql"
)

func GetMigrations(t *testing.T, ctx context.Context) []string {
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		t.Error("Unable to get caller information")
	}

	basePath := filepath.Dir(b)
	migrationsAbsDir := filepath.Join(basePath, migrationsDir)

	files, err := filepath.Glob(filepath.Join(migrationsAbsDir, "*"+migrationsSuffix))
	if err != nil {
		t.Errorf("Failed to read migrations: %v", err)
	}

	return files
}

func CreatePostgresContainer(t *testing.T, ctx context.Context) (*postgres.PostgresContainer, error) {
	t.Helper()

	migrations := GetMigrations(t, ctx)
	pgContainer, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithInitScripts(migrations...),
		postgres.WithDatabase("goyurback"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Errorf("Could not start postgres container: %v", err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("Could not terminate postgres container: %v", err)
		}
	})

	return pgContainer, nil
}
