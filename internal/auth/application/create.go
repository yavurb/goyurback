package application

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
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

	sha512Hash := sha512.New()
	sha512Hash.Write([]byte(keyString))
	sha := sha512Hash.Sum(nil)
	hashedApikey := hex.EncodeToString(sha)

	keySalt := strings.Split(keyString, ".")[0]
	hashedApikeyWithSalt := strings.Join([]string{keySalt, hashedApikey}, ".")

	apiKey := &domain.APIKeyCreate{
		Key:  hashedApikeyWithSalt,
		Name: name,
	}

	createdKey, err := uc.repository.CreateAPIKey(ctx, apiKey)
	if err != nil {
		return nil, err
	}

	keyString = strings.Join([]string{prefix, keyString}, "_")
	createdKey.Key = keyString

	return createdKey, nil
}
