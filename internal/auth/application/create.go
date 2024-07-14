package application

import (
	"context"
	"strings"

	"github.com/yavurb/goyurback/internal/auth/domain"
	"github.com/yavurb/goyurback/internal/pgk/ids"
)

const prefix = "sk"

func (uc *apiKeyUsecase) CreateAPIKey(ctx context.Context, name string) (*domain.APIKey, error) {
	keyString, err := ids.NewAPIKey()
	if err != nil {
		return nil, err
	}

	keyString = strings.Join([]string{prefix, keyString}, "_")

	// TODO: Hash the key before saving, and return the unhashed key

	apiKey := &domain.APIKeyCreate{
		Key:  keyString,
		Name: name,
	}

	createdKey, err := uc.repository.CreateAPIKey(ctx, apiKey)
	if err != nil {
		return nil, err
	}

	return createdKey, nil
}
