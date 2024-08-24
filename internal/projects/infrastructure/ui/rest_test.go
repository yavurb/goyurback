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

func TestGetProject(t *testing.T) {
	e := echo.New()

	t.Run("it should get a project", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/projects/:id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("pr_12345")

		uc := &mocks.MockProjectsUsecase{}
		h := NewProjectsRouter(e, uc)

		want := map[string]any{
			"id":            "pr_12345",
			"name":          "test",
			"description":   "Some project description",
			"tags":          []string{"tag1", "tag2", "tag3"},
			"thumbnail_url": "https://example.com/image.jpg",
			"website_url":   "https://example.com",
			"live":          true,
			"created_at":    time.Now().UTC().Format(time.RFC3339Nano),
			"updated_at":    time.Now().UTC().Format(time.RFC3339Nano),
		}

		uc.GetFn = func(ctx context.Context, id string) (*domain.Project, error) {
			if id != "pr_12345" {
				return nil, domain.ErrProjectNotFound
			}

			createdAt, _ := time.Parse(time.RFC3339, want["created_at"].(string))
			updatedAt, _ := time.Parse(time.RFC3339, want["updated_at"].(string))

			return &domain.Project{
				ID:           1,
				PublicID:     "pr_12345",
				Name:         "test",
				Description:  "Some project description",
				ThumbnailURL: "https://example.com/image.jpg",
				WebsiteURL:   "https://example.com",
				Live:         true,
				Tags:         []string{"tag1", "tag2", "tag3"},
				PostID:       1,
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
			}, nil
		}

		err := h.getProject(c)
		if err != nil {
			t.Errorf("getProject() error = %v, want no error", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("getProject() status = %v, want %v", rec.Code, http.StatusOK)
		}

		got := make(map[string]any)
		if err = json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Errorf("Error unmarshalling response: %s", err)
		}

		if !testhelpers.CompareMaps(want, got) {
			t.Errorf("getProject() mismatch:\n%s", cmp.Diff(want, got))
		}
	})

	t.Run("it should return a bad request error", func(t *testing.T) {
		t.Skip("Test not implemented yet")

		req := httptest.NewRequest(http.MethodGet, "/projects/:id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("/projects/:id")

		uc := &mocks.MockProjectsUsecase{}
		h := NewProjectsRouter(e, uc)

		uc.GetFn = func(ctx context.Context, id string) (*domain.Project, error) {
			return nil, domain.ErrProjectNotFound
		}

		err := h.getProject(c)
		if err == nil {
			t.Error("Got nil getting project, want error")
		}

		if !errors.Is(err, echo.ErrBadRequest) {
			t.Errorf("Error getting project = %v, want %v", err, echo.ErrBadRequest)
		}
	})

	t.Run("it should return a not found error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/projects/:id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("/projects/:id")
		c.SetParamNames("id")
		c.SetParamValues("pr_12345")

		uc := &mocks.MockProjectsUsecase{}
		h := NewProjectsRouter(e, uc)

		uc.GetFn = func(ctx context.Context, id string) (*domain.Project, error) {
			return nil, domain.ErrProjectNotFound
		}

		err := h.getProject(c)
		if err == nil {
			t.Errorf("getProject() error = %v, want an error", err)
		}

		if !errors.Is(err, echo.ErrNotFound) {
			t.Errorf("getProject() error = %v, want %v", err, echo.ErrNotFound)
		}

		if !strings.Contains(err.Error(), "Project not found") {
			t.Errorf("getProject() body = %v, want %v", rec.Body.String(), "Project not found")
		}
	})
}

func TestGetProjects(t *testing.T) {
	e := echo.New()

	t.Run("it should return the projects", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/projects", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("/projects")

		uc := &mocks.MockProjectsUsecase{}
		h := NewProjectsRouter(e, uc)

		want := map[string][]map[string]any{
			"data": {
				{
					"id":            "pr_12345",
					"name":          "test",
					"description":   "Some project description",
					"tags":          []string{"tag1", "tag2", "tag3"},
					"thumbnail_url": "https://example.com/image.jpg",
					"website_url":   "https://example.com",
					"live":          true,
					"created_at":    time.Now().UTC().Format(time.RFC3339Nano),
					"updated_at":    time.Now().UTC().Format(time.RFC3339Nano),
				},
				{
					"id":            "pr_12346",
					"name":          "test 2",
					"description":   "Some second project description",
					"tags":          []string{"tag4", "tag5", "tag6"},
					"thumbnail_url": "https://example.com/image2.jpg",
					"website_url":   "https://exampletwo.com",
					"live":          true,
					"created_at":    time.Now().UTC().Format(time.RFC3339Nano),
					"updated_at":    time.Now().UTC().Format(time.RFC3339Nano),
				},
			},
		}

		uc.GetProjectsFn = func(ctx context.Context) ([]*domain.Project, error) {
			createdAt, _ := time.Parse(time.RFC3339, want["data"][0]["created_at"].(string))
			updatedAt, _ := time.Parse(time.RFC3339, want["data"][0]["updated_at"].(string))
			createdAt2, _ := time.Parse(time.RFC3339, want["data"][1]["created_at"].(string))
			updatedAt2, _ := time.Parse(time.RFC3339, want["data"][1]["updated_at"].(string))

			return []*domain.Project{
				{
					ID:           1,
					PublicID:     "pr_12345",
					Name:         "test",
					Description:  "Some project description",
					ThumbnailURL: "https://example.com/image.jpg",
					WebsiteURL:   "https://example.com",
					Live:         true,
					Tags:         []string{"tag1", "tag2", "tag3"},
					PostID:       1,
					CreatedAt:    createdAt,
					UpdatedAt:    updatedAt,
				},
				{
					ID:           2,
					PublicID:     "pr_12346",
					Name:         "test 2",
					Description:  "Some second project description",
					ThumbnailURL: "https://example.com/image2.jpg",
					WebsiteURL:   "https://exampletwo.com",
					Live:         true,
					Tags:         []string{"tag4", "tag5", "tag6"},
					PostID:       2,
					CreatedAt:    createdAt2,
					UpdatedAt:    updatedAt2,
				},
			}, nil
		}

		err := h.getProjects(c)
		if err != nil {
			t.Errorf("getProjects() error = %v, want no error", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("getProjects() status = %v, want %v", rec.Code, http.StatusOK)
		}

		got := make(map[string]any)
		if err = json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Errorf("Error unmarshalling response: %s", err)
		}

		if !testhelpers.CompareMaps(want, got) {
			t.Errorf("getProject() mismatch:\n%s", cmp.Diff(want, got))
		}
	})

	t.Run("it should return a internal server error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/projects", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("/projects")

		uc := &mocks.MockProjectsUsecase{}
		h := NewProjectsRouter(e, uc)

		uc.GetProjectsFn = func(ctx context.Context) ([]*domain.Project, error) {
			return nil, errors.New("DB Error")
		}

		err := h.getProjects(c)
		if err == nil {
			t.Error("Got nil gettings projects, want error")
		}

		if !errors.Is(err, echo.ErrInternalServerError) {
			t.Errorf("getProject() error = %v, want %v", err, echo.ErrNotFound)
		}
	})
}
