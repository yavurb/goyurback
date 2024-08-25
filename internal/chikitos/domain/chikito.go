package domain

import (
	"time"

	"github.com/google/go-cmp/cmp"
)

type Chikito struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PublicID    string
	URL         string
	Description string
	ID          int32
}

func (c Chikito) Compare(c2 Chikito) bool {
	c.CreatedAt = c2.CreatedAt
	c.UpdatedAt = c2.UpdatedAt

	return cmp.Equal(c, c2)
}

type ChikitoCreate struct {
	PublicID    string
	URL         string
	Description string
}
