package publicid

import (
	"strings"

	nanoid "github.com/matoous/go-nanoid/v2"
)

const (
	alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
	lenght   = 12
)

func New(prefix string) (string, error) {
	id, err := nanoid.Generate(alphabet, lenght)

	if err != nil {
		return "", err
	}

	return strings.Join([]string{prefix, id}, "_"), nil
}
