package application

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
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
	}

	repo := &mocks.MockPostsRepository{
		CreatePostFn: func(ctx context.Context, post *domain.PostCreate) (*domain.Post, error) {
			return &domain.Post{
				ID:          1,
				PublicID:    post.PublicID,
				Title:       post.Title,
				Author:      post.Author,
				Slug:        post.Slug,
				Status:      domain.Draft,
				Description: post.Description,
				Content:     post.Content,
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
			}, nil
		},
	}

	uc := NewPostUsecase(repo)
	ctx := context.Background()

	got, err := uc.Create(ctx, want.Title, want.Author, want.Slug, want.Description, want.Content)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	rgx := regexp.MustCompile(`po_[a-zA-Z0-9]{5}`)
	if !rgx.MatchString(got.PublicID) {
		t.Errorf("Expected PublicID to match the regex %s, got: %s", rgx.String(), got.PublicID)
	}

	want.PublicID = got.PublicID
	if !want.Compare(*got) {
		t.Errorf("Mismatch creating post (-want,+got):\n%s", cmp.Diff(want, got))
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
