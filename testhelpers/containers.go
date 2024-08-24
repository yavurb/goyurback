package testhelpers

import (
	"context"
	"log"
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

type PostgresContainer struct {
	Container  *postgres.PostgresContainer
	ConnString string
}

func CreatePostgresContainer(t *testing.T, ctx context.Context) (*PostgresContainer, error) {
	t.Helper()

	pgContainer, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
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

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Errorf("Could not get connection string: %v", err)
	}

	ApplyMigrations(t, ctx, connStr)

	return &PostgresContainer{
		Container:  pgContainer,
		ConnString: connStr,
	}, nil
}
