package ids

import (
	"strings"

	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/yavurb/goyurback/internal/pgk/rand"
)

const (
	alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
	length   = 12
)

var generateRandomString = rand.GenerateRandomString

func NewPublicID(prefix string) (string, error) {
	id, err := nanoid.Generate(alphabet, length)
	if err != nil {
		return "", err
	}

	return strings.Join([]string{prefix, id}, "_"), nil
}

func NewAPIKey() (string, error) {
	apiKey, err := generateRandomString(64)
	if err != nil {
		return "", err
	}

	salt, err := generateRandomString(32)
	if err != nil {
		return "", err
	}

	// Remove the special characters from the generated string
	apiKey = strings.ReplaceAll(apiKey, "-", "")
	apiKey = strings.ReplaceAll(apiKey, "_", "")
	apiKey = strings.ReplaceAll(apiKey, "=", "")

	salt = strings.ReplaceAll(salt, "=", "")
	salt = strings.ReplaceAll(salt, "-", "")
	salt = strings.ReplaceAll(salt, "_", "")

	apiKey = strings.Join([]string{salt, apiKey}, ".")

	return apiKey, nil
}
