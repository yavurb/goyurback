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

func TestGetPosts(t *testing.T) {
	e := echo.New()
	e.Validator = mods.NewAppValidator()

	t.Run("it should return the posts", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/posts", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		uc := &mocks.MockPostsUsecase{}

		want := map[string][]map[string]any{
			"data": {
				{
					"id":           "po_12345",
					"title":        "My test post",
					"author":       "Roy",
					"slug":         "my-test-post",
					"status":       "published",
					"description":  "Some post description",
					"content":      "Some post content",
					"created_at":   time.Now().UTC().Format(time.RFC3339),
					"updated_at":   time.Now().UTC().Format(time.RFC3339),
					"published_at": time.Now().UTC().Format(time.RFC3339),
				},
				{
					"id":           "po_12346",
					"title":        "My test post 2",
					"author":       "Roy 2",
					"slug":         "my-test-post 2",
					"status":       "published",
					"description":  "Some post description 2",
					"content":      "Some post content 2",
					"created_at":   time.Now().UTC().Format(time.RFC3339),
					"updated_at":   time.Now().UTC().Format(time.RFC3339),
					"published_at": time.Now().UTC().Format(time.RFC3339),
				},
			},
		}

		uc.GetPostsFn = func(ctx context.Context) ([]*domain.Post, error) {
			createdAt, _ := time.Parse(time.RFC3339, want["data"][0]["created_at"].(string))
			updatedAt, _ := time.Parse(time.RFC3339, want["data"][0]["updated_at"].(string))
			publishedAt, _ := time.Parse(time.RFC3339, want["data"][0]["published_at"].(string))
			createdAt2, _ := time.Parse(time.RFC3339, want["data"][1]["created_at"].(string))
			updatedAt2, _ := time.Parse(time.RFC3339, want["data"][1]["updated_at"].(string))
			publishedAt2, _ := time.Parse(time.RFC3339, want["data"][1]["published_at"].(string))

			return []*domain.Post{
				{
					ID:          1,
					PublicID:    "po_12345",
					Title:       "My test post",
					Author:      "Roy",
					Slug:        "my-test-post",
					Status:      domain.Published,
					Description: "Some post description",
					Content:     "Some post content",
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
					PublishedAt: publishedAt,
				},
				{
					ID:          2,
					PublicID:    "po_12346",
					Title:       "My test post 2",
					Author:      "Roy 2",
					Slug:        "my-test-post 2",
					Status:      domain.Published,
					Description: "Some post description 2",
					Content:     "Some post content 2",
					CreatedAt:   createdAt2,
					UpdatedAt:   updatedAt2,
					PublishedAt: publishedAt2,
				},
			}, nil
		}

		h := NewPostsRouter(e, uc)

		err := h.getPosts(c)
		if err != nil {
			t.Errorf("Expected no errors getting posts. Got: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected http status to be a 200 StatusOK. Got: %d", rec.Code)
		}

		got := make(map[string]any)
		if err = json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Errorf("Error unmarshalling response: %s", err)
		}

		if !testhelpers.CompareMaps(want, got) {
			t.Errorf("Mismatch gettings posts:\n%s", cmp.Diff(want, got))
		}
	})

	t.Run("it should return a internal server error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/posts", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		uc := &mocks.MockPostsUsecase{}
		h := NewPostsRouter(e, uc)

		uc.GetPostsFn = func(ctx context.Context) ([]*domain.Post, error) {
			return nil, errors.New("DB Error")
		}

		err := h.getPosts(c)
		if err == nil {
			t.Error("Got nil gettings posts, want error")
		}

		if !errors.Is(err, echo.ErrInternalServerError) {
			t.Errorf("Expected error to be a 500 Internal Server Error. Got: %v", err)
		}
	})
}
