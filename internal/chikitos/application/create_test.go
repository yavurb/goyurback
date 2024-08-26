package application

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/yavurb/goyurback/internal/chikitos/application/mocks"
	"github.com/yavurb/goyurback/internal/chikitos/domain"
)

func TestCreate(t *testing.T) {
	t.Run("it should create a chikito", func(t *testing.T) {
		want := &domain.Chikito{
			URL:         "https://example.com/my_long_url",
			Description: "My long URL description",
			ID:          1,
		}
		repo := &mocks.MockChikitosRepository{}

		repo.CreateChikitoFn = func(ctx context.Context, chikito *domain.ChikitoCreate) (*domain.Chikito, error) {
			return &domain.Chikito{
				ID:          1,
				PublicID:    chikito.PublicID,
				URL:         chikito.URL,
				Description: chikito.Description,
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
			}, nil
		}

		uc := NewChikitoUsecase(repo)

		got, err := uc.Create(context.Background(), want.URL, want.Description)
		if err != nil {
			t.Errorf("Expected no error creating chikito, got %v", err)
		}

		rgx := regexp.MustCompile(`ch_[a-zA-Z0-9]{5}`)
		if !rgx.MatchString(got.PublicID) {
			t.Errorf("Expected PublicID to match the regex %s, got: %s", rgx.String(), got.PublicID)
		}

		want.PublicID = got.PublicID
		if !want.Compare(*got) {
			t.Errorf("Mismatch creating chikito. (-want,+got):\n%v", cmp.Diff(want, got))
		}
	})

	t.Run("it should return an error when creating a chikito", func(t *testing.T) {
		repo := &mocks.MockChikitosRepository{}

		repo.CreateChikitoFn = func(ctx context.Context, chikito *domain.ChikitoCreate) (*domain.Chikito, error) {
			return nil, errors.New("DB Error")
		}

		uc := NewChikitoUsecase(repo)

		_, err := uc.Create(context.Background(), "https://example.com", "some description")
		if err == nil {
			t.Error("Expected error creating chikito, got nil")
		}
	})
}
