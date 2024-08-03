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

func TestCreatePost(t *testing.T) {
	want := &domain.Post{
		ID:          1,
		PublicID:    "someid",
		Title:       "Some Post",
		Author:      "Some Author",
		Slug:        "some-post",
		Status:      "draft",
		Description: "Some Description",
		Content:     "Some Content",
		PublishedAt: time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	repo := &mocks.MockPostsRepository{
		CreatePostFn: func(ctx context.Context, post *domain.PostCreate) (*domain.Post, error) {
			return &domain.Post{
				ID:          1,
				PublicID:    "someid",
				Title:       post.Title,
				Author:      post.Author,
				Slug:        post.Slug,
				Status:      "draft",
				Description: post.Description,
				Content:     post.Content,
				PublishedAt: want.PublishedAt,
				CreatedAt:   want.CreatedAt,
				UpdatedAt:   want.UpdatedAt,
			}, nil
		},
	}

	uc := NewPostUsecase(repo)
	ctx := context.Background()

	post, err := uc.Create(ctx, want.Title, want.Author, want.Slug, want.Description, want.Content)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !reflect.DeepEqual(post, want) {
		t.Errorf("Expected post to be %v, got: %v", want, post)
	}
}

func TestCreatePostWithDBError(t *testing.T) {
	want := "unable to create post"
	repo := &mocks.MockPostsRepository{
		CreatePostFn: func(ctx context.Context, post *domain.PostCreate) (*domain.Post, error) {
			return nil, errors.New("DB Error")
		},
	}

	uc := NewPostUsecase(repo)
	ctx := context.Background()

	_, err := uc.Create(ctx, "Some post", "Royner Perez", "Some Slug", "Some Description", "Some content")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	if err.Error() != want {
		t.Errorf("Expected error to be %v, got: %v", want, err)
	}
}
