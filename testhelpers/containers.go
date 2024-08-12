package testhelpers

import (
	"context"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func GetMigrations() []string {
	// FIXME: This is a hacky way to get the path to the migrations directory.
	migrations := []string{
		"20240323001046_initial.up.sql",
		"20240714014750_apikeys.up.sql",
		"20240715004121_variable_apikey_length.up.sql",
	}

	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)

	for i, migration := range migrations {
		migrations[i] = filepath.Join(basePath, "../migrations", migration)
	}

	return migrations
}

func CreatePostgresContainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	migrations := GetMigrations()
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
		log.Fatalf("Could not start postgres container: %s", err)
	}

	return pgContainer, nil
}
