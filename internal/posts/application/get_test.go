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

func TestGet(t *testing.T) {
	want := &domain.Post{
		ID:          1,
		PublicID:    "public-id",
		Title:       "some title",
		Author:      "some author",
		Slug:        "some-slug",
		Status:      domain.Published,
		Description: "some description",
		Content:     "<h1>some content</h1>",
		PublishedAt: time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	repo := &mocks.MockPostsRepository{
		GetPostFn: func(ctx context.Context, id string) (*domain.Post, error) {
			if id != want.PublicID {
				return nil, errors.New("post not found")
			}

			return want, nil
		},
	}

	uc := NewPostUsecase(repo)

	post, err := uc.Get(context.Background(), want.PublicID)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !reflect.DeepEqual(post, want) {
		t.Errorf("Expected post to be %v, got: %v", want, post)
	}
}

func TestGetPostNotFound(t *testing.T) {
	repo := &mocks.MockPostsRepository{
		GetPostFn: func(ctx context.Context, id string) (*domain.Post, error) {
			return nil, errors.New("DB error")
		},
	}

	uc := NewPostUsecase(repo)

	_, err := uc.Get(context.Background(), "non-existing-id")

	if !errors.Is(err, domain.ErrPostNotFound) {
		t.Errorf("Expected error to be %v, got: %v", domain.ErrPostNotFound, err)
	}
}
