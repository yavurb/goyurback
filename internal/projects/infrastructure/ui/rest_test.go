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
	"github.com/yavurb/goyurback/internal/projects/domain"
	"github.com/yavurb/goyurback/internal/projects/infrastructure/ui/mocks"
	"github.com/yavurb/goyurback/testhelpers"
)

func TestCreateProject(t *testing.T) {
	e := echo.New()

	want := map[string]any{
		"id":            "pr_12345",
		"name":          "test",
		"description":   "Some project description",
		"tags":          []string{"tag1", "tag2", "tag3"},
		"thumbnail_url": "https://example.com/image.jpg",
		"website_url":   "https://example.com",
		"live":          true,
		"created_at":    time.Now().UTC().Format(time.RFC3339),
		"updated_at":    time.Now().UTC().Format(time.RFC3339),
	}

	t.Run("it should create a project", func(t *testing.T) {
		jsonBytes, err := json.Marshal(want)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/projects", strings.NewReader(string(jsonBytes)))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		uc := &mocks.MockProjectsUsecase{
			CreateFn: func(ctx context.Context, name, description, thumbnailURL, websiteURL string, live bool, tags []string, postId int32) (*domain.Project, error) {
				createdAt, _ := time.Parse(time.RFC3339, want["created_at"].(string))
				updatedAt, _ := time.Parse(time.RFC3339, want["updated_at"].(string))
				return &domain.Project{
					ID:           1,
					PublicID:     "pr_12345",
					Name:         name,
					Description:  description,
					ThumbnailURL: thumbnailURL,
					WebsiteURL:   websiteURL,
					Live:         live,
					Tags:         tags,
					PostID:       postId,
					CreatedAt:    createdAt,
					UpdatedAt:    updatedAt,
				}, nil
			},
		}
		h := NewProjectsRouter(e, uc)

		err = h.createProject(c)
		if err != nil {
			t.Errorf("createProject() error = %v, want no error", err)
		}

		if rec.Code != http.StatusCreated {
			t.Errorf("createProject() status = %v, want %v", rec.Code, http.StatusCreated)
		}

		got := make(map[string]any)
		if err = json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Errorf("Error unmarshalling response: %s", err)
		}

		if !testhelpers.CompareMaps(want, got) {
			t.Errorf("createProject() mismatch:\n%s", cmp.Diff(want, got))
		}
	})

	t.Run("it should return an internal server error", func(t *testing.T) {
		project := map[string]any{
			"name":          "test",
			"description":   "Some project description",
			"tags":          []string{"tag1", "tag2", "tag3"},
			"thumbnail_url": "https://example.com/image.jpg",
			"website_url":   "https://example.com",
			"live":          true,
		}

		jsonBytes, err := json.Marshal(project)
		if err != nil {
			t.Fatal(err)
		}

		payload := strings.NewReader(string(jsonBytes))
		req := httptest.NewRequest(http.MethodPost, "/projects", payload)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		uc := &mocks.MockProjectsUsecase{CreateFn: func(
			ctx context.Context,
			name, description, thumbnailURL, websiteURL string,
			live bool,
			tags []string,
			postId int32,
		) (*domain.Project, error) {
			return nil, errors.New("Unknown usecase error")
		}}
		h := NewProjectsRouter(e, uc)

		err = h.createProject(c)
		if err == nil {
			t.Errorf("createProject() error = %v, want an error", err)
		}

		if !errors.Is(err, echo.ErrInternalServerError) {
			t.Errorf("createProject() error = %v, want %v", err, echo.ErrInternalServerError)
		}

		if !strings.Contains(err.Error(), "Internal server error") {
			t.Errorf("createProject() body = %v, want %v", rec.Body.String(), "Invalid request body")
		}
	})

	t.Run("it should return an error", func(t *testing.T) {
		project := map[string]any{
			"name":          "test",
			"description":   "Some project description",
			"tags":          []string{"tag1", "tag2", "tag3"},
			"thumbnail_url": "https://example.com/image.jpg",
			"website_url":   "https://example.com",
			"live":          "true",
		}

		jsonBytes, err := json.Marshal(project)
		if err != nil {
			t.Fatal(err)
		}

		payload := strings.NewReader(string(jsonBytes))
		req := httptest.NewRequest(http.MethodPost, "/projects", payload)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		uc := &mocks.MockProjectsUsecase{}
		h := NewProjectsRouter(e, uc)

		err = h.createProject(c)
		if err == nil {
			t.Errorf("createProject() error = %v, want an error", err)
		}

		if !errors.Is(err, echo.ErrUnprocessableEntity) {
			t.Errorf("createProject() error = %v, want %v", err, echo.ErrUnprocessableEntity)
		}

		if !strings.Contains(err.Error(), "Invalid request body") {
			t.Errorf("createProject() body = %v, want %v", rec.Body.String(), "Invalid request body")
		}
	})
}
