package rand

import (
	"encoding/base64"
	"testing"
)

func isBase64(got string) bool {
	_, err := base64.URLEncoding.DecodeString(got)
	return err == nil
}

func TestGenerateRandomBytes(t *testing.T) {
	t.Run("it should generate random bytes", func(t *testing.T) {
		n := 32

		got, err := GenerateRandomBytes(n)
		if err != nil {
			t.Errorf("GenerateRandomBytes() error = %v, want nil", err)
		}

		if len(got) != n {
			t.Errorf("GenerateRandomBytes() len = %v, want %v", len(got), n)
		}
	})
}

func TestGenerateRandomString(t *testing.T) {
	t.Run("it should generate random string", func(t *testing.T) {
		n := 32

		got, err := GenerateRandomString(n)
		if err != nil {
			t.Errorf("GenerateRandomString() error = %v, want nil", err)
		}

		if !isBase64(got) {
			t.Errorf("GenerateRandomString() = %v, want base64 encoded string", got)
		}

		secondGot, _ := GenerateRandomString(n)
		if got == secondGot {
			t.Errorf("GenerateRandomString() = %v, want different", got)
		}
	})
}
