package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yavurb/goyurback/internal/posts/domain"
	"github.com/yavurb/goyurback/testhelpers"
)

func TestCreatePost(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := testhelpers.CreatePostgresContainer(t, ctx)
	if err != nil {
		t.Fatal(err)
	}

	conn, err := pgxpool.New(ctx, pgContainer.ConnString)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { conn.Close() })

	repo := NewRepo(conn)

	t.Run("it should create a post", func(t *testing.T) {
		testhelpers.CleanDatabase(t, ctx, pgContainer.ConnString)

		want := &domain.Post{
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			PublicID:    "po_18892",
			Title:       "My Post",
			Author:      "Roy",
			Slug:        "my-post",
			Status:      domain.Draft,
			Description: "my post description",
			Content:     "# My Post\n\nThis is my post content.",
			ID:          1,
		}

		got, err := repo.CreatePost(ctx, &domain.PostCreate{
			PublicID:    "po_18892",
			Title:       "My Post",
			Author:      "Roy",
			Slug:        "my-post",
			Description: "my post description",
			Content:     "# My Post\n\nThis is my post content.",
		})
		if err != nil {
			t.Errorf("Got error creating post, want no error: %v", err)
		}

		if !want.Compare(*got) {
			t.Errorf("Mismatch creating project (-want,+got):\n%s", cmp.Diff(want, got))
		}
	})

	t.Run("it should return an error if a db error occurs", func(t *testing.T) {
		testhelpers.DeleteDatabase(t, ctx, pgContainer.ConnString)

		_, err := repo.CreatePost(ctx, &domain.PostCreate{
			PublicID:    "po_18892",
			Title:       "My Post",
			Author:      "Roy",
			Slug:        "my-post",
			Description: "my post description",
			Content:     "# My Post\n\nThis is my post content.",
		})

		if err == nil {
			t.Error("Got no error, want error")
		}
	})
}
