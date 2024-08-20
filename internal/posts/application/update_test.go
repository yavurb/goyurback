package application

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/yavurb/goyurback/internal/posts/application/mocks"
	"github.com/yavurb/goyurback/internal/posts/domain"
)

type postUpdate struct {
	Status      *domain.Status
	Title       *string
	Author      *string
	Slug        *string
	Description *string
	Content     *string
	ID          string
}

var post = domain.Post{
	ID:          1,
	PublicID:    "pk_1",
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

// Helper function to create string pointers
func pointer[T any](s T) *T {
	return &s
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		want     *domain.Post
		toUpdate postUpdate
	}{
		{func() *domain.Post { temp := post; temp.Status = domain.Archived; return &temp }(), postUpdate{pointer(domain.Archived), nil, nil, nil, nil, nil, "pk_1"}},
		{func() *domain.Post { temp := post; temp.Title = "title"; return &temp }(), postUpdate{nil, pointer("title"), nil, nil, nil, nil, "pk_1"}},
		{func() *domain.Post { temp := post; temp.Author = "new author"; return &temp }(), postUpdate{nil, nil, pointer("new author"), nil, nil, nil, "pk_1"}},
		{func() *domain.Post { temp := post; temp.Slug = "new-slug"; return &temp }(), postUpdate{nil, nil, nil, pointer("new-slug"), nil, nil, "pk_1"}},
		{func() *domain.Post { temp := post; temp.Description = "new description"; return &temp }(), postUpdate{nil, nil, nil, nil, pointer("new description"), nil, "pk_1"}},
		{func() *domain.Post { temp := post; temp.Content = "<h1>new content</h1>"; return &temp }(), postUpdate{nil, nil, nil, nil, nil, pointer("<h1>new content</h1>"), "pk_1"}},
	}

	repo := &mocks.MockPostsRepository{
		GetPostFn: func(ctx context.Context, id string) (*domain.Post, error) {
			postCopy := post
			return &postCopy, nil
		},
		UpdatePostFn: func(context context.Context, p *domain.Post) (*domain.Post, error) {
			return p, nil
		},
	}

	uc := NewPostUsecase(repo)

	for _, test := range tests {
		testName := fmt.Sprintf("it should update field %s", structToString(test.toUpdate))
		t.Run(testName, func(t *testing.T) {
			postUpdated, err := uc.Update(
				context.Background(),
				test.toUpdate.ID,
				test.toUpdate.Title,
				test.toUpdate.Author,
				test.toUpdate.Slug,
				test.toUpdate.Description,
				test.toUpdate.Content,
				test.toUpdate.Status,
			)
			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}

			if postUpdated == nil {
				t.Errorf("Expected post to be %v, got: %v", test.want, postUpdated)
			}

			if !reflect.DeepEqual(postUpdated, test.want) {
				t.Errorf("Expected post to be %v, got: %v", test.want, postUpdated)
			}
		})
	}
}

func structToString(post postUpdate) string {
	var parts []string

	if post.Status != nil {
		parts = append(parts, "status")
	}

	if post.Title != nil {
		parts = append(parts, "title")
	}

	if post.Author != nil {
		parts = append(parts, "author")
	}

	if post.Slug != nil {
		parts = append(parts, "slug")
	}

	if post.Description != nil {
		parts = append(parts, "description")
	}

	if post.Content != nil {
		parts = append(parts, "content")
	}

	return strings.Join(parts, ":")
}
