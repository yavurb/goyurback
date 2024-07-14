package application

import (
	"context"

	"github.com/yavurb/goyurback/internal/auth/domain"
	"github.com/yavurb/goyurback/internal/pgk/ids"
)

func (uc *apiKeyUsecase) CreateAPIKey(ctx context.Context, keyName, key string) (*domain.APIKey, error) {
	keyString, err := ids.NewAPIKey()
	if err != nil {
		return nil, err
	}

	// TODO: Hash the key before saving

	apiKey := &domain.APIKeyCreate{
		Key:     keyString,
		KeyName: keyName,
	}

	createdKey, err := uc.repository.CreateAPIKey(ctx, apiKey)
	if err != nil {
		return nil, err
	}

	return createdKey, nil
}
