package ui

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/labstack/echo/v4"
	"github.com/yavurb/goyurback/internal/app/mods"
	"github.com/yavurb/goyurback/internal/posts/domain"
	"github.com/yavurb/goyurback/internal/posts/infrastructure/ui/mocks"
	"github.com/yavurb/goyurback/testhelpers"
)

func TestCreatePost(t *testing.T) {
	e := echo.New()
	e.Validator = mods.NewAppValidator()

	t.Run("it should create a post", func(t *testing.T) {
		postIn := map[string]any{
			"title":       "My test post",
			"author":      "Roy",
			"slug":        "my-test-post",
			"description": "Some post description",
			"content":     "Some post content",
		}
		want := map[string]any{
			"id":           "po_12345",
			"title":        "My test post",
			"author":       "Roy",
			"slug":         "my-test-post",
			"status":       "draft",
			"description":  "Some post description",
			"content":      "Some post content",
			"created_at":   time.Now().UTC().Format(time.RFC3339),
			"updated_at":   time.Now().UTC().Format(time.RFC3339),
			"published_at": new(time.Time).Format(time.RFC3339),
		}

		jsonBytes, err := json.Marshal(postIn)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(string(jsonBytes)))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		uc := &mocks.MockPostsUsecase{}
		uc.CreateFn = func(ctx context.Context, title, author, slug, description, content string) (*domain.Post, error) {
			createdAt, _ := time.Parse(time.RFC3339, want["created_at"].(string))
			updatedAt, _ := time.Parse(time.RFC3339, want["updated_at"].(string))

			return &domain.Post{
				ID:          1,
				PublicID:    "po_12345",
				Title:       title,
				Author:      author,
				Slug:        slug,
				Status:      domain.Draft,
				Description: description,
				Content:     content,
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
			}, nil
		}

		h := NewPostsRouter(e, uc)

		err = h.createPost(c)
		if err != nil {
			t.Errorf("Expected no error creating post. Got: %v", err)
		}

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected http status code to be %d. Got %d", http.StatusCreated, rec.Code)
		}

		got := make(map[string]any)
		if err = json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Errorf("Error unmarshalling response: %v", err)
		}

		if !testhelpers.CompareMaps(want, got) {
			t.Errorf("Mismatch creating post:\n%s", cmp.Diff(want, got))
		}
	})

	t.Run("it should return an internal server error", func(t *testing.T) {
		post := map[string]any{
			"title":       "My test post",
			"author":      "Roy",
			"slug":        "my-test-post",
			"status":      domain.Draft,
			"description": "Some post description",
			"content":     "Some post content",
		}

		jsonBytes, err := json.Marshal(post)
		if err != nil {
			t.Fatal(err)
		}

		payload := strings.NewReader(string(jsonBytes))
		req := httptest.NewRequest(http.MethodPost, "/posts", payload)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		uc := &mocks.MockPostsUsecase{}
		uc.CreateFn = func(ctx context.Context, title, author, slug, description, content string) (*domain.Post, error) {
			return nil, errors.New("Unknown usecase error")
		}

		h := NewPostsRouter(e, uc)

		err = h.createPost(c)
		if err == nil {
			t.Error("Expected error creating post. Got nil")
		}

		if !errors.Is(err, echo.ErrInternalServerError) {
			t.Errorf("Expected error to be a 500 (ErrInternalServerError). Got: %v", err)
		}
	})

	t.Run("it should return a validation error", func(t *testing.T) {
		post := map[string]any{
			"title":       "My test post",
			"author":      1,
			"slug":        "my-test-post",
			"status":      domain.Draft,
			"description": "Some post description",
			"content":     "Some post content",
		}

		jsonBytes, err := json.Marshal(post)
		if err != nil {
			t.Fatal(err)
		}

		payload := strings.NewReader(string(jsonBytes))
		req := httptest.NewRequest(http.MethodPost, "/posts", payload)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		uc := &mocks.MockPostsUsecase{}
		h := NewPostsRouter(e, uc)

		err = h.createPost(c)
		if err == nil {
			t.Error("Expected error creating post. Got nil")
		}

		if !errors.Is(err, echo.ErrUnprocessableEntity) {
			t.Errorf("Expected error to be a 422 (ErrUnprocessableEntity). Got: %v", err)
		}
	})
}

func TestGetPost(t *testing.T) {
	e := echo.New()
	e.Validator = mods.NewAppValidator()

	t.Run("it should get a post", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/posts/:id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("po_12345")

		uc := &mocks.MockPostsUsecase{}

		want := map[string]any{
			"id":           "po_12345",
			"title":        "My test post",
			"author":       "Roy",
			"slug":         "my-test-post",
			"status":       domain.Draft,
			"description":  "Some post description",
			"content":      "Some post content",
			"created_at":   time.Now().UTC().Format(time.RFC3339),
			"updated_at":   time.Now().UTC().Format(time.RFC3339),
			"published_at": new(time.Time).Format(time.RFC3339),
		}

		uc.GetFn = func(ctx context.Context, id string) (*domain.Post, error) {
			if id != "po_12345" {
				return nil, domain.ErrPostNotFound
			}

			createdAt, _ := time.Parse(time.RFC3339, want["created_at"].(string))
			updatedAt, _ := time.Parse(time.RFC3339, want["updated_at"].(string))

			return &domain.Post{
				ID:          1,
				PublicID:    "po_12345",
				Title:       "My test post",
				Author:      "Roy",
				Slug:        "my-test-post",
				Status:      domain.Draft,
				Description: "Some post description",
				Content:     "Some post content",
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
			}, nil
		}

		h := NewPostsRouter(e, uc)

		err := h.getPost(c)
		if err != nil {
			t.Errorf("Expected no errors getting post. Got: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status code to be a 200 (StatusOK). Got: %d", rec.Code)
		}

		got := make(map[string]any)
		if err = json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Errorf("Error unmarshalling response: %s", err)
		}

		if !testhelpers.CompareMaps(want, got) {
			t.Errorf("Mismatch:\n%s", cmp.Diff(want, got))
		}
	})

	t.Run("it should return a bad unprocessable entity", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/posts/:id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("/posts/:id")

		uc := &mocks.MockPostsUsecase{}
		h := NewPostsRouter(e, uc)

		err := h.getPost(c)
		if err == nil {
			t.Error("Expected error getting post. Got nil")
		}

		if !errors.Is(err, echo.ErrUnprocessableEntity) {
			t.Errorf("Expected request error to be a 422 (ErrUnprocessableEntity). Got: %v", err)
		}
	})

	t.Run("it should return a not found error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/posts/:id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("/posts/:id")
		c.SetParamNames("id")
		c.SetParamValues("po_12345")

		uc := &mocks.MockPostsUsecase{}
		uc.GetFn = func(ctx context.Context, id string) (*domain.Post, error) {
			return nil, domain.ErrPostNotFound
		}

		h := NewPostsRouter(e, uc)

		err := h.getPost(c)
		if err == nil {
			t.Error("Expected error getting post. Got nil")
		}

		if !errors.Is(err, echo.ErrNotFound) {
			t.Errorf("Expected request error to be a 404 (ErrNotFound). Got: %v", err)
		}
	})
}

// func TestGetPosts(t *testing.T) {
// 	e := echo.New()
//
// 	t.Run("it should return the posts", func(t *testing.T) {
// 		req := httptest.NewRequest(http.MethodGet, "/posts", nil)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)
//
// 		c.SetPath("/posts")
//
// 		uc := &mocks.MockPostsUsecase{}
// 		h := NewPostsRouter(e, uc)
//
// 		want := map[string][]map[string]any{
// 			"data": {
// 				{
// 					"id":            "pr_12345",
// 					"name":          "test",
// 					"description":   "Some post description",
// 					"tags":          []string{"tag1", "tag2", "tag3"},
// 					"thumbnail_url": "https://example.com/image.jpg",
// 					"website_url":   "https://example.com",
// 					"live":          true,
// 					"created_at":    time.Now().UTC().Format(time.RFC3339Nano),
// 					"updated_at":    time.Now().UTC().Format(time.RFC3339Nano),
// 				},
// 				{
// 					"id":            "pr_12346",
// 					"name":          "test 2",
// 					"description":   "Some second post description",
// 					"tags":          []string{"tag4", "tag5", "tag6"},
// 					"thumbnail_url": "https://example.com/image2.jpg",
// 					"website_url":   "https://exampletwo.com",
// 					"live":          true,
// 					"created_at":    time.Now().UTC().Format(time.RFC3339Nano),
// 					"updated_at":    time.Now().UTC().Format(time.RFC3339Nano),
// 				},
// 			},
// 		}
//
// 		uc.GetPostsFn = func(ctx context.Context) ([]*domain.Post, error) {
// 			createdAt, _ := time.Parse(time.RFC3339, want["data"][0]["created_at"].(string))
// 			updatedAt, _ := time.Parse(time.RFC3339, want["data"][0]["updated_at"].(string))
// 			createdAt2, _ := time.Parse(time.RFC3339, want["data"][1]["created_at"].(string))
// 			updatedAt2, _ := time.Parse(time.RFC3339, want["data"][1]["updated_at"].(string))
//
// 			return []*domain.Post{
// 				{
// 					ID:           1,
// 					PublicID:     "pr_12345",
// 					Name:         "test",
// 					Description:  "Some post description",
// 					ThumbnailURL: "https://example.com/image.jpg",
// 					WebsiteURL:   "https://example.com",
// 					Live:         true,
// 					Tags:         []string{"tag1", "tag2", "tag3"},
// 					PostID:       1,
// 					CreatedAt:    createdAt,
// 					UpdatedAt:    updatedAt,
// 				},
// 				{
// 					ID:           2,
// 					PublicID:     "pr_12346",
// 					Name:         "test 2",
// 					Description:  "Some second post description",
// 					ThumbnailURL: "https://example.com/image2.jpg",
// 					WebsiteURL:   "https://exampletwo.com",
// 					Live:         true,
// 					Tags:         []string{"tag4", "tag5", "tag6"},
// 					PostID:       2,
// 					CreatedAt:    createdAt2,
// 					UpdatedAt:    updatedAt2,
// 				},
// 			}, nil
// 		}
//
// 		err := h.getPosts(c)
// 		if err != nil {
// 			t.Errorf("getPosts() error = %v, want no error", err)
// 		}
//
// 		if rec.Code != http.StatusOK {
// 			t.Errorf("getPosts() status = %v, want %v", rec.Code, http.StatusOK)
// 		}
//
// 		got := make(map[string]any)
// 		if err = json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
// 			t.Errorf("Error unmarshalling response: %s", err)
// 		}
//
// 		if !testhelpers.CompareMaps(want, got) {
// 			t.Errorf("getPost() mismatch:\n%s", cmp.Diff(want, got))
// 		}
// 	})
//
// 	t.Run("it should return a internal server error", func(t *testing.T) {
// 		req := httptest.NewRequest(http.MethodGet, "/posts", nil)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)
//
// 		c.SetPath("/posts")
//
// 		uc := &mocks.MockPostsUsecase{}
// 		h := NewPostsRouter(e, uc)
//
// 		uc.GetPostsFn = func(ctx context.Context) ([]*domain.Post, error) {
// 			return nil, errors.New("DB Error")
// 		}
//
// 		err := h.getPosts(c)
// 		if err == nil {
// 			t.Error("Got nil gettings posts, want error")
// 		}
//
// 		if !errors.Is(err, echo.ErrInternalServerError) {
// 			t.Errorf("getPost() error = %v, want %v", err, echo.ErrNotFound)
// 		}
// 	})
// }
