package testhelpers

import (
	"encoding/json"

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
	json.Unmarshal(aBytes, &aT)
	json.Unmarshal(bBytes, &bT)

	return cmp.Equal(aT, bT)
}
