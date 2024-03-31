package ui

import "time"

type CreateIn struct {
	URL         string `json:"url" validate:"required, url"`
	Description string `json:"description" validate:"required"`
}

type CreateOut struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
