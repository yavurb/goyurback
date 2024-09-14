package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yavurb/goyurback/internal/chikitos/domain"
	"github.com/yavurb/goyurback/testhelpers"
)

func TestCreateChikito(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := testhelpers.CreatePostgresContainer(t, ctx)
	if err != nil {
		t.Fatalf("Error creating postgres container: %v", err)
	}

	t.Run("It should create a chikito", func(t *testing.T) {
		testhelpers.CleanDatabase(t, ctx, pgContainer.ConnString)

		conn, err := pgxpool.New(ctx, pgContainer.ConnString)
		if err != nil {
			t.Fatalf("Error creating pgxpool: %v", err)
		}

		t.Cleanup(func() { conn.Close() })

		repo := NewRepo(conn)

		want := &domain.Chikito{
			ID:          1,
			PublicID:    "ch_12345",
			URL:         "https://example.com/my_long_url",
			Description: "My long URL description",
		}

		got, err := repo.CreateChikito(ctx, &domain.ChikitoCreate{
			PublicID:    "ch_12345",
			URL:         "https://example.com/my_long_url",
			Description: "My long URL description",
		})
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if !want.Compare(*got) {
			t.Errorf("Mismatch creating chikito. (-want,+got):\n%s", cmp.Diff(want, got))
		}
	})

	t.Run("It should return an error if the public id already exists", func(t *testing.T) {
		testhelpers.CleanDatabase(t, ctx, pgContainer.ConnString)

		conn, err := pgxpool.New(ctx, pgContainer.ConnString)
		if err != nil {
			t.Fatalf("Error creating pgxpool: %v", err)
		}

		t.Cleanup(func() { conn.Close() })

		repo := NewRepo(conn)

		_, err = repo.CreateChikito(ctx, &domain.ChikitoCreate{
			PublicID:    "ch_12345",
			URL:         "https://example.com/my_long_url",
			Description: "My long URL description",
		})
		if err != nil {
			t.Errorf("Expected no error creating first chikito, got: %v", err)
		}

		_, err = repo.CreateChikito(ctx, &domain.ChikitoCreate{
			PublicID:    "ch_12345",
			URL:         "https://example.com/my_long_url",
			Description: "My long URL description",
		})
		if err == nil {
			t.Errorf("Expected error creating second chikito, got nil")
		}

		if !errors.Is(err, domain.ErrPublicIDAlreadyExists) {
			t.Errorf("Expected error creating second chikito to be ErrPublicIDAlreadyExists, got: %v", err)
		}
	})

	t.Run("It should return an error if there is a db error", func(t *testing.T) {
		testhelpers.CleanDatabase(t, ctx, pgContainer.ConnString)
		testhelpers.DeleteDatabase(t, ctx, pgContainer.ConnString)

		conn, err := pgxpool.New(ctx, pgContainer.ConnString)
		if err != nil {
			t.Fatalf("Error creating pgxpool: %v", err)
		}

		t.Cleanup(func() { conn.Close() })

		repo := NewRepo(conn)

		_, err = repo.CreateChikito(ctx, &domain.ChikitoCreate{
			PublicID:    "ch_12345",
			URL:         "https://example.com/my_long_url",
			Description: "My long URL description",
		})
		if err == nil {
			t.Errorf("Expected error creating second chikito, got nil")
		}

		pgErr := new(pgconn.PgError)
		if !errors.As(err, &pgErr) {
			t.Errorf("Expected error to be a pgx.PgError, got: %v", err)
		}
	})
}

func TestGetChikito(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := testhelpers.CreatePostgresContainer(t, ctx)
	if err != nil {
		t.Fatalf("Error creating postgres container: %v", err)
	}

	t.Run("it should get a chikito", func(t *testing.T) {
		testhelpers.CleanDatabase(t, ctx, pgContainer.ConnString)

		conn, err := pgxpool.New(ctx, pgContainer.ConnString)
		if err != nil {
			t.Fatalf("Error creating pgxpool: %v", err)
		}

		t.Cleanup(func() { conn.Close() })

		repo := NewRepo(conn)

		_, err = repo.CreateChikito(ctx, &domain.ChikitoCreate{
			PublicID:    "ch_12345",
			URL:         "https://example.com/my_long_url",
			Description: "My long URL description",
		})
		if err != nil {
			t.Errorf("Expected no error creating chikito, got: %v", err)
		}

		want := &domain.Chikito{
			ID:          1,
			PublicID:    "ch_12345",
			URL:         "https://example.com/my_long_url",
			Description: "My long URL description",
		}

		got, err := repo.GetChikito(ctx, "ch_12345")
		if err != nil {
			t.Errorf("Expected no error getting chikito. Got: %v", err)
		}

		if !want.Compare(*got) {
			t.Errorf("Mismatch getting chikito. (-want,+got):\n%s", cmp.Diff(want, got))
		}
	})

	t.Run("It should return a not found error", func(t *testing.T) {
		testhelpers.CleanDatabase(t, ctx, pgContainer.ConnString)

		conn, err := pgxpool.New(ctx, pgContainer.ConnString)
		if err != nil {
			t.Fatalf("Error creating pgxpool: %v", err)
		}

		t.Cleanup(func() { conn.Close() })

		repo := NewRepo(conn)

		_, err = repo.GetChikito(ctx, "ch_12345")
		if err == nil {
			t.Errorf("Expected error getting chikito. Got nil")
		}

		if !errors.Is(err, domain.ErrChikitoNotFound) {
			t.Errorf("Expected error to be ErrChikitoNotFound, got: %v", err)
		}
	})

	t.Run("It should return an error if there is a db error", func(t *testing.T) {
		testhelpers.CleanDatabase(t, ctx, pgContainer.ConnString)
		testhelpers.DeleteDatabase(t, ctx, pgContainer.ConnString)

		conn, err := pgxpool.New(ctx, pgContainer.ConnString)
		if err != nil {
			t.Fatalf("Error creating pgxpool: %v", err)
		}

		t.Cleanup(func() { conn.Close() })

		repo := NewRepo(conn)

		_, err = repo.GetChikito(ctx, "ch_12345")
		if err == nil {
			t.Errorf("Expected error creating second chikito, got nil")
		}

		pgErr := new(pgconn.PgError)
		if !errors.As(err, &pgErr) {
			t.Errorf("Expected error to be a pgx.PgError, got: %v", err)
		}
	})
}
