package ui

import "time"

type APIKeyIn struct {
	Name string `json:"name" validate:"required,min=5,max=64"`
}

type APIKeyOut struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Key       string    `json:"key"`
	CreatedAt time.Time `json:"created_at"`
}
