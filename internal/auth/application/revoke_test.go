package application

import (
	"context"
	"testing"

	"github.com/yavurb/goyurback/internal/auth/application/mocks"
)

func TestRevokeAPIKey(t *testing.T) {
	repo := &mocks.MockAuthRepository{
		RevokeAPIKeyFn: func(ctx context.Context, publicID string) error {
			return nil
		},
	}

	uc := NewAPIKeyUsecase(repo)

	ctx := context.Background()
	err := uc.RevokeAPIKey(ctx, "random-id")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}
