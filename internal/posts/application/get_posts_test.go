package application

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/yavurb/goyurback/internal/posts/application/mocks"
	"github.com/yavurb/goyurback/internal/posts/domain"
)

func TestGetPosts(t *testing.T) {
	want := []*domain.Post{
		{
			ID:          1,
			PublicID:    "public-id-1",
			Title:       "some title",
			Author:      "some author",
			Slug:        "some-slug-1",
			Status:      domain.Published,
			Description: "some description",
			Content:     "<h1>some content</h1>",
			PublishedAt: time.Now(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			PublicID:    "public-id-2",
			Title:       "some new title",
			Author:      "some new author",
			Slug:        "some-slug-2",
			Status:      domain.Archived,
			Description: "some new description",
			Content:     "<h1>some new content</h1>",
			PublishedAt: time.Now(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	repo := &mocks.MockPostsRepository{
		GetPostsFn: func(ctx context.Context) ([]*domain.Post, error) {
			return want, nil
		},
	}

	uc := NewPostUsecase(repo)

	posts, err := uc.GetPosts(context.Background())
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !reflect.DeepEqual(posts, want) {
		t.Errorf("Expected posts to be %v, got: %v", want, posts)
	}
}

func TestGetPostsError(t *testing.T) {
	repo := &mocks.MockPostsRepository{
		GetPostsFn: func(ctx context.Context) ([]*domain.Post, error) {
			return nil, errors.New("DB error")
		},
	}

	uc := NewPostUsecase(repo)

	_, err := uc.GetPosts(context.Background())

	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	if err.Error() != "DB error" {
		t.Errorf("Expected error to be 'DB error', got: %v", err)
	}
}
