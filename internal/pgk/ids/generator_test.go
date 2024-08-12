package ids

import (
	"errors"
	"regexp"
	"strings"
	"testing"
)

func TestNewPublicID(t *testing.T) {
	t.Run("it should generate random public id", func(t *testing.T) {
		prefix := "prefix"
		got, err := NewPublicID(prefix)
		if err != nil {
			t.Errorf("GenerateID() error = %v, want nil", err)
		}
		if !strings.HasPrefix(got, prefix) {
			t.Errorf("GenerateID() = %v, want to has prefix %v", got, prefix)
		}

		secondGot, _ := NewPublicID(prefix)
		if got == secondGot {
			t.Errorf("GenerateID() = %v, want different", got)
		}
	})
}

func TestNewAPIKey(t *testing.T) {
	t.Run("it should generate random api key", func(t *testing.T) {
		got, err := NewAPIKey()
		if err != nil {
			t.Errorf("GenerateAPIKey() error = %v, want nil", err)
		}

		pattern := `^[a-zA-Z0-9]+\.[a-zA-Z0-9]+$`
		if !regexp.MustCompile(pattern).MatchString(got) {
			t.Errorf("GenerateAPIKey() = %v, want to satisfy patter: %v", got, pattern)
		}
		secondGot, _ := NewAPIKey()
		if got == secondGot {
			t.Errorf("GenerateAPIKey() = %v, want different", got)
		}
	})

	t.Run("it should return an error if key generation fails", func(t *testing.T) {
		errorMessage := "key generation error"
		generateRandomString = func(_ int) (string, error) { return "", errors.New(errorMessage) }

		_, err := NewAPIKey()
		if err == nil {
			t.Error("GenerateAPIKey() error is nil, want error")
		}
		if err.Error() != errorMessage {
			t.Errorf(`GenerateAPIKey() error = %v, want "%v"`, err, errorMessage)
		}
	})

	t.Run("it should return an error if salt generation fails", func(t *testing.T) {
		count := 0
		errorMessage := "salt generation error"
		generateRandomString = func(_ int) (string, error) {
			if count < 1 {
				count++
				return "key", nil
			} else {
				return "", errors.New(errorMessage)
			}
		}

		_, err := NewAPIKey()
		if err == nil {
			t.Error("GenerateAPIKey() error is nil, want error")
		}
		if err.Error() != errorMessage {
			t.Errorf(`GenerateAPIKey() error = %v, want "%v"`, err, errorMessage)
		}
	})
}
