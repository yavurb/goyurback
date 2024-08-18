package ui

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/yavurb/goyurback/internal/projects/application"
	"github.com/yavurb/goyurback/internal/projects/infrastructure/repository"
	"github.com/yavurb/goyurback/testhelpers"
)

func TestCreateProject(t *testing.T) {
	e := echo.New()

	ctx := context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(t, ctx)
	if err != nil {
		t.Errorf("Error creating container: %s", err)
	}

	connpool, err := pgxpool.New(ctx, pgContainer.ConnString)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	t.Cleanup(func() { connpool.Close() })

	project := map[string]any{
		"name":          "test",
		"description":   "Some project description",
		"tags":          []string{"tag1", "tag2", "tag3"},
		"thumbnail_url": "https://example.com/image.jpg",
		"website_url":   "https://example.com",
		"live":          true,
	}
	rec := httptest.NewRecorder()

	t.Run("it should create a project", func(t *testing.T) {
		testhelpers.CleanDatabase(t, ctx, pgContainer.ConnString)

		want := map[string]any{
			"id":            1,
			"name":          "test",
			"description":   "Some project description",
			"tags":          []string{"tag1", "tag2", "tag3"},
			"thumbnail_url": "https://example.com/image.jpg",
			"website_url":   "https://example.com",
			"live":          true,
		}

		jsonString, err := json.Marshal(project)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/projects", strings.NewReader(string(jsonString)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)

		repo := repository.NewRepo(connpool)
		uc := application.NewProjectUsecase(repo)
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
		want["id"] = got["id"]
		want["created_at"] = got["created_at"]
		want["updated_at"] = got["updated_at"]

		wantStr, _ := json.Marshal(want)
		gotStr, _ := json.Marshal(got)

		if string(wantStr) != string(gotStr) {
			t.Errorf("createProject() = %s, want %s", string(gotStr), string(wantStr))
		}
	})
}
