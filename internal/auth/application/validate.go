package application

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/yavurb/goyurback/internal/auth/domain"
)

func (uc *apiKeyUsecase) ValidateAPIKey(ctx context.Context, key string) (bool, error) {
	match, err := regexp.MatchString("^[a-z]{2}_[a-zA-Z0-9]+\\.[a-zA-Z0-9]+$", key)
	if err != nil {
		log.Printf("Error matching api key: %v\n", err)

		return false, err
	}
	if !match {
		return false, nil
	}

	apiKey := strings.Split(key, "_")[1] // remove prefix
	apiKeySalt := strings.Split(apiKey, ".")[0]

	sha512Hash := sha512.New()
	sha512Hash.Write([]byte(apiKey))
	sha := sha512Hash.Sum(nil)
	hashedApikey := hex.EncodeToString(sha)

	_, err = uc.repository.GetAPIKeyByValue(
		ctx,
		strings.Join([]string{apiKeySalt, hashedApikey}, "."),
	)
	if err != nil {
		if errors.Is(err, domain.ErrAPIKeyNotFound) {
			return false, nil
		}

		log.Printf("Error getting api key by value: %v\n", err)

		return false, err
	}

	return true, nil
}
