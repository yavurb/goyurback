package repository

import (
	"context"
	"errors"
	"slices"
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

func TestGetPost(t *testing.T) {
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

	t.Run("it should get a post", func(t *testing.T) {
		testhelpers.CleanDatabase(t, ctx, pgContainer.ConnString)

		_, _ = repo.CreatePost(ctx, &domain.PostCreate{
			PublicID:    "po_18892",
			Title:       "My Post",
			Author:      "Roy",
			Slug:        "my-post",
			Description: "my post description",
			Content:     "# My Post\n\nThis is my post content.",
		})

		want := &domain.Post{
			PublicID:    "po_18892",
			Title:       "My Post",
			Author:      "Roy",
			Slug:        "my-post",
			Status:      domain.Draft,
			Description: "my post description",
			Content:     "# My Post\n\nThis is my post content.",
			ID:          1,
		}

		got, err := repo.GetPost(ctx, "po_18892")
		if err != nil {
			t.Errorf("Got error getting post, want no error: %v", err)
		}

		if !want.Compare(*got) {
			t.Errorf("Mismatch getting post (-want,+got):\n%s", cmp.Diff(want, got))
		}
	})

	// NOTE: This test has its own database because a an error with cached data types in postgres
	t.Run("it should return a not found error", func(t *testing.T) {
		conn, err := pgxpool.New(ctx, pgContainer.ConnString)
		if err != nil {
			t.Fatal(err)
		}

		t.Cleanup(func() { conn.Close() })

		repo := NewRepo(conn)

		_, err = repo.GetPost(ctx, "po_18810")
		if err == nil {
			t.Errorf("Got nil getting post, want error")
		}

		if !errors.Is(err, domain.ErrPostNotFound) {
			t.Errorf("Error mismatch gettings post. Got %v, want %v", err, domain.ErrPostNotFound)
		}
	})
}

func TestGetPosts(t *testing.T) {
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

	t.Run("it should get all posts", func(t *testing.T) {
		testhelpers.CleanDatabase(t, ctx, pgContainer.ConnString)

		want := []*domain.Post{
			{
				PublicID:    "po_18892",
				Title:       "My Post",
				Author:      "Roy",
				Slug:        "my-post",
				Status:      domain.Published,
				Description: "my post description",
				Content:     "# My Post\n\nThis is my post content.",
				ID:          1,
			},
			{
				PublicID:    "po_18893",
				Title:       "My Post 2",
				Author:      "Roy",
				Slug:        "my-post-2",
				Status:      domain.Published,
				Description: "my post description 2",
				Content:     "# My Post\n\nThis is my post content 2.",
				ID:          2,
			},
		}
		for _, w := range want {
			postCreated, err := repo.CreatePost(ctx, &domain.PostCreate{
				PublicID:    w.PublicID,
				Title:       w.Title,
				Author:      w.Author,
				Slug:        w.Slug,
				Description: w.Description,
				Content:     w.Content,
			})
			if err != nil {
				t.Errorf("Got error creating post, want no error: %v", err)
			}

			postCreated.Status = domain.Published
			postCreated.PublishedAt = time.Now()

			_, err = repo.UpdatePost(ctx, postCreated)
			if err != nil {
				t.Errorf("Got error updating post, want no error: %v", err)
			}
		}

		got, err := repo.GetPosts(ctx)
		if err != nil {
			t.Errorf("Got error getting posts, want no error: %v", err)
		}

		slices.Reverse(want)

		for i, w := range want {
			if !w.Compare(*got[i]) {
				t.Errorf("Mismatch getting posts (-want,+got):\n%s", cmp.Diff(w, got[i]))
			}
		}
	})
}
