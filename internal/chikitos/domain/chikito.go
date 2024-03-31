package domain

import "time"

type Chikito struct {
	ID          int32
	PublicID    string
	URL         string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
