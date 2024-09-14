package application

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/yavurb/goyurback/internal/auth/application/mocks"
	"github.com/yavurb/goyurback/internal/auth/domain"
)

func TestCreateAPIKey(t *testing.T) {
	repo := &mocks.MockAuthRepository{
		CreateAPIKeyFn: func(ctx context.Context, apiKey *domain.APIKeyCreate) (*domain.APIKey, error) {
			return &domain.APIKey{
				ID:       1,
				PublicID: "test",
				Key:      apiKey.Key,
				Name:     apiKey.Name,
				Revoked:  false,
			}, nil
		},
	}

	uc := NewAPIKeyUsecase(repo)

	ctx := context.Background()

	apikey, err := uc.CreateAPIKey(ctx, "test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !strings.HasPrefix(apikey.Key, "sk") {
		t.Errorf("Expected key to start with 'sk', got %s", apikey.Key)
	}

	if apikey.PublicID != "test" {
		t.Errorf("Expected publicID to be 'test', got %s", apikey.PublicID)
	}

	match, err := regexp.MatchString("^[a-z]{2}_[a-zA-Z0-9]+\\.[a-zA-Z0-9]+$", apikey.Key)
	if err != nil {
		t.Errorf("Error matching api key: %v", err)
	}

	if !match {
		t.Errorf("Expected key to match regex pattern xx_xxx+.xxxxxx+, got %s", apikey.Key)
	}
}
