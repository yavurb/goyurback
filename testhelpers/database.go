package testhelpers

import (
	"context"
	"log"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func ApplyMigrations(t *testing.T, ctx context.Context, connStr string) {
	t.Helper()

	_, b, _, ok := runtime.Caller(0)
	if !ok {
		t.Error("Unable to get caller information")
	}

	basePath := filepath.Dir(b)
	migrationsAbsDir := filepath.Join("file:///", basePath, "../migrations")

	m, err := migrate.New(migrationsAbsDir, connStr)
	if err != nil {
		t.Fatal(err)
	}

	err = m.Up()
	if err != nil {
		t.Fatalf("Error migrating up database: %v", err)
	}
}

func CleanDatabase(t *testing.T, ctx context.Context, connStr string) {
	t.Helper()

	t.Cleanup(func() {
		_, b, _, ok := runtime.Caller(0)
		if !ok {
			t.Error("Unable to get caller information")
		}

		basePath := filepath.Dir(b)
		migrationsAbsDir := filepath.Join("file:///", basePath, "../migrations")

		m, err := migrate.New(migrationsAbsDir, connStr)
		if err != nil {
			t.Fatal(err)
		}

		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			t.Fatalf("Error dropping database: %v", err)
		}

		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Printf("Error migrating up: %v", err)
			log.Fatal(err)
		}
	})
}

func DeleteDatabase(t *testing.T, ctx context.Context, connStr string) {
	t.Helper()

	_, b, _, ok := runtime.Caller(0)
	if !ok {
		t.Error("Unable to get caller information")
	}

	basePath := filepath.Dir(b)
	migrationsAbsDir := filepath.Join("file:///", basePath, "../migrations")

	m, err := migrate.New(migrationsAbsDir, connStr)
	if err != nil {
		t.Fatal(err)
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("Error dropping database: %v", err)
	}

	t.Cleanup(func() {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Printf("Error migrating up: %v", err)
			log.Fatal(err)
		}
	})
}
