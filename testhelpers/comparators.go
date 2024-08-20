package testhelpers

import (
	"encoding/json"
	"log"

	"github.com/google/go-cmp/cmp"
)

func CompareMaps(a, b map[string]any) bool {
	if len(a) != len(b) {
		return false
	}

	aBytes, _ := json.Marshal(a)
	bBytes, _ := json.Marshal(b)
	aT := make(map[string]any)
	bT := make(map[string]any)

	err := json.Unmarshal(aBytes, &aT)
	if err != nil {
		log.Fatal("Error unmarshalling aBytes")
	}

	err = json.Unmarshal(bBytes, &bT)
	if err != nil {
		log.Fatal("Error unmarshalling bBytes")
	}

	return cmp.Equal(aT, bT)
}
