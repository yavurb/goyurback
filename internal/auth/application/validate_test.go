package application

import (
	"context"
	"testing"
	"time"

	"github.com/yavurb/goyurback/internal/auth/application/mocks"
	"github.com/yavurb/goyurback/internal/auth/domain"
)

const (
	apiKey     string = "sk_testsalt.testkey"
	apiKeyHash string = "testsalt.a2eff62ef47f9d1313d3755a43d42883a010fcaf1faebc43207defb25fa58621789bbd33292923efea4f4d44b43a3ab4bb9d8b2d6f0577867bac8a4874242ba4"
)

func TestValidate(t *testing.T) {
	t.Run("it should return true", func(t *testing.T) {
		repo := &mocks.MockAuthRepository{
			GetAPIKeyByValueFn: func(ctx context.Context, apiKey string) (*domain.APIKey, error) {
				if apiKey != apiKeyHash {
					return nil, domain.ErrAPIKeyNotFound
				}

				return &domain.APIKey{
					CreatedAt: time.Now().UTC(),
					UpdatedAt: time.Now().UTC(),
					Name:      "testing_key",
					Key:       apiKeyHash,
					PublicID:  "ak_1qjrblb8pm90",
					ID:        1,
					Revoked:   false,
				}, nil
			},
		}

		uc := NewAPIKeyUsecase(repo)
		ctx := context.Background()

		valid, err := uc.ValidateAPIKey(ctx, apiKey)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if !valid {
			t.Errorf("Expected valid to be true, got: %v", valid)
		}
	})

	t.Run("it should return false", func(t *testing.T) {
		repo := &mocks.MockAuthRepository{
			GetAPIKeyByValueFn: func(ctx context.Context, apiKey string) (*domain.APIKey, error) {
				if apiKey != apiKeyHash {
					return nil, domain.ErrAPIKeyNotFound
				}

				return &domain.APIKey{
					CreatedAt: time.Now().UTC(),
					UpdatedAt: time.Now().UTC(),
					Name:      "testing_key",
					Key:       apiKeyHash,
					PublicID:  "ak_1qjrblb8pm90",
					ID:        1,
					Revoked:   false,
				}, nil
			},
		}

		uc := NewAPIKeyUsecase(repo)
		ctx := context.Background()

		valid, err := uc.ValidateAPIKey(ctx, "sk_invalid.key")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if valid == true {
			t.Errorf("Expected valid to be true, got: %v", valid)
		}
	})

	t.Run("it should return false if the key does not follow the key format sk_salt.hash", func(t *testing.T) {
		repo := &mocks.MockAuthRepository{}
		uc := NewAPIKeyUsecase(repo)
		ctx := context.Background()

		valid, err := uc.ValidateAPIKey(ctx, "invalid-key")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if valid == true {
			t.Errorf("Expected valid to be true, got: %v", valid)
		}
	})
}
