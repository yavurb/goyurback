package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/yavurb/goyurback/internal/chikitos/application/mocks"
	"github.com/yavurb/goyurback/internal/chikitos/domain"
)

func TestGet(t *testing.T) {
	t.Run("it should get a chikito", func(t *testing.T) {
		want := &domain.Chikito{
			ID:          1,
			PublicID:    "ch_12345",
			URL:         "https://example.com/my_long_url",
			Description: "My long URL description",
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}

		repo := &mocks.MockChikitosRepository{}
		repo.GetChikitoFn = func(ctx context.Context, id string) (*domain.Chikito, error) {
			return want, nil
		}
		uc := NewChikitoUsecase(repo)

		got, err := uc.Get(context.Background(), "ch_12345")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if !want.Compare(*got) {
			t.Errorf("Mismatch getting chikito. (-want,+got):\n%v", cmp.Diff(want, got))
		}
	})

	t.Run("it should return a not found error", func(t *testing.T) {
		repo := &mocks.MockChikitosRepository{}
		repo.GetChikitoFn = func(ctx context.Context, id string) (*domain.Chikito, error) {
			return nil, domain.ErrChikitoNotFound
		}
		uc := NewChikitoUsecase(repo)

		_, err := uc.Get(context.Background(), "ch_12345")
		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		if !errors.Is(err, domain.ErrChikitoNotFound) {
			t.Errorf("Expected ErrChikitoNotFound error, got: %v", err)
		}
	})
}
