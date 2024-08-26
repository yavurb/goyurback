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
	"github.com/yavurb/goyurback/internal/chikitos/domain"
	"github.com/yavurb/goyurback/internal/chikitos/infrastructure/ui/mocks"
	"github.com/yavurb/goyurback/testhelpers"
)

func TestCreateChikito(t *testing.T) {
	e := echo.New()
	e.Validator = mods.NewAppValidator()

	t.Run("It should create a chikito", func(t *testing.T) {
		chikitoIn := map[string]any{
			"url":         "https://example.com",
			"description": "Some random description",
		}
		want := map[string]any{
			"id":          "ch_12345",
			"url":         "https://example.com",
			"description": "Some random description",
			"created_at":  time.Now().UTC().Format(time.RFC3339),
			"updated_at":  time.Now().UTC().Format(time.RFC3339),
		}

		jsonBytes, err := json.Marshal(chikitoIn)
		if err != nil {
			t.Fatal(err)
		}

		jsonString := string(jsonBytes)

		req := httptest.NewRequest(http.MethodPost, "/chikitos", strings.NewReader(jsonString))
		rec := httptest.NewRecorder()

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		c := e.NewContext(req, rec)
		uc := &mocks.MockChikitosUsecase{}
		uc.CreateFn = func(ctx context.Context, url, description string) (*domain.Chikito, error) {
			createdAt, _ := time.Parse(time.RFC3339, want["created_at"].(string))
			updatedAt, _ := time.Parse(time.RFC3339, want["updated_at"].(string))

			return &domain.Chikito{
				ID:          1,
				PublicID:    "ch_12345",
				URL:         url,
				Description: description,
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
			}, nil
		}

		h := NewChikitosRouter(e, uc)

		err = h.create(c)
		if err != nil {
			t.Errorf("Expected no error creating chikito, got %v", err)
		}

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, rec.Code)
		}

		got := make(map[string]any)
		if err = json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Errorf("Error unmarshalling response: %s", err)
		}

		if !testhelpers.CompareMaps(want, got) {
			t.Errorf("Mismatch creating chikito:\n%s", cmp.Diff(want, got))
		}
	})

	t.Run("It should return a 422 error", func(t *testing.T) {
		chikitoIn := map[string]any{
			"description": "Some random description",
		}

		jsonBytes, err := json.Marshal(chikitoIn)
		if err != nil {
			t.Fatal(err)
		}

		jsonString := string(jsonBytes)

		req := httptest.NewRequest(http.MethodPost, "/chikitos", strings.NewReader(jsonString))
		rec := httptest.NewRecorder()

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		c := e.NewContext(req, rec)
		uc := &mocks.MockChikitosUsecase{}
		uc.CreateFn = func(ctx context.Context, url, description string) (*domain.Chikito, error) {
			return &domain.Chikito{}, nil
		}

		h := NewChikitosRouter(e, uc)

		err = h.create(c)
		if err == nil {
			t.Errorf("Expected error creating chikito, got nil")
		}

		if !errors.Is(err, echo.ErrUnprocessableEntity) {
			t.Errorf("Expected error to be echo.ErrUnprocessableEntity, got %v", err)
		}
	})

	t.Run("It should return a internal server error", func(t *testing.T) {
		chikitoIn := map[string]any{
			"url":         "https://example.com",
			"description": "Some random description",
		}

		jsonBytes, err := json.Marshal(chikitoIn)
		if err != nil {
			t.Fatal(err)
		}

		jsonString := string(jsonBytes)

		req := httptest.NewRequest(http.MethodPost, "/chikitos", strings.NewReader(jsonString))
		rec := httptest.NewRecorder()

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		c := e.NewContext(req, rec)
		uc := &mocks.MockChikitosUsecase{}
		uc.CreateFn = func(ctx context.Context, url, description string) (*domain.Chikito, error) {
			return nil, domain.ErrPublicIDAlreadyExists
		}

		h := NewChikitosRouter(e, uc)

		err = h.create(c)
		if err == nil {
			t.Errorf("Expected error creating chikito, got nil")
		}

		if !errors.Is(err, echo.ErrInternalServerError) {
			t.Errorf("Expected error to be echo.ErrInternalServerError, got %v", err)
		}
	})
}

func TestGetChikito(t *testing.T) {
	e := echo.New()
	e.Validator = mods.NewAppValidator()

	t.Run("It should redirect to the chikito's URL", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/chikitos/:id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("ch_12345")

		uc := &mocks.MockChikitosUsecase{}
		uc.GetFn = func(ctx context.Context, id string) (*domain.Chikito, error) {
			if id != "ch_12345" {
				return nil, domain.ErrChikitoNotFound
			}

			return &domain.Chikito{
				ID:          1,
				PublicID:    "ch_12345",
				URL:         "https://example.com",
				Description: "Some random description",
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
			}, nil
		}

		h := NewChikitosRouter(e, uc)

		err := h.get(c)
		if err != nil {
			t.Errorf("Expected no error getting chikito. Got: %v", err)
		}

		if rec.Code != http.StatusPermanentRedirect {
			t.Errorf("Expected response code to be a 308 (StatusPermanentRedirect). Got: %d", rec.Code)
		}

		if rec.Result().Header.Get("Location") != "https://example.com" {
			t.Errorf("Expected location to be 'https://example.com'. Got: %s", rec.Result().Header.Get("Location"))
		}
	})

	t.Run("It should return a not found error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/chikitos/:id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("ch_12345")

		uc := &mocks.MockChikitosUsecase{}
		uc.GetFn = func(ctx context.Context, id string) (*domain.Chikito, error) {
			return nil, domain.ErrChikitoNotFound
		}

		h := NewChikitosRouter(e, uc)

		err := h.get(c)
		if err == nil {
			t.Error("Expected error getting chikito. Got nil")
		}

		if !errors.Is(err, echo.ErrNotFound) {
			t.Errorf("Expected error to be a 404 (ErrNotFound). Got: %v", err)
		}
	})

	t.Run("It should return a unprocessable entity error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/chikitos/:id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("/:id")
		c.SetParamNames("id")

		uc := &mocks.MockChikitosUsecase{}
		uc.GetFn = func(ctx context.Context, id string) (*domain.Chikito, error) {
			return nil, domain.ErrChikitoNotFound
		}

		h := NewChikitosRouter(e, uc)

		err := h.get(c)
		if err == nil {
			t.Error("Expected error getting chikito. Got nil")
		}

		if !errors.Is(err, echo.ErrUnprocessableEntity) {
			t.Errorf("Expected error to be a 422 (ErrUnprocessableEntity). Got: %v", err)
		}
	})
}
