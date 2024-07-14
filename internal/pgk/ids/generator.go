package ids

import (
	"strings"

	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/yavurb/goyurback/internal/pgk/rand"
)

const (
	alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
	lenght   = 12
)

func NewPublicID(prefix string) (string, error) {
	id, err := nanoid.Generate(alphabet, lenght)
	if err != nil {
		return "", err
	}

	return strings.Join([]string{prefix, id}, "_"), nil
}

func NewAPIKey() (string, error) {
	apiKey, err := rand.GenerateRandomString(64)
	if err != nil {
		return "", err
	}

	// Remove the special characters from the generated string
	apiKey = strings.ReplaceAll(apiKey, "-", "")
	apiKey = strings.ReplaceAll(apiKey, "_", "")
	apiKey = strings.ReplaceAll(apiKey, "=", "")

	return apiKey, nil
}
