package ui

import (
	"time"

	"github.com/yavurb/goyurback/internal/posts/domain"
)

type PostIn struct {
	Title       string `json:"title" validate:"required,min=5,max=128"`
	Author      string `json:"author" validate:"required,min=3,max=64"`
	Slug        string `json:"slug" validate:"required"`
	Description string `json:"description" validate:"required,min=5,max=255"`
	Content     string `json:"content" validate:"required,min=10"`
}

type PostOut struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Author      string        `json:"author"`
	Slug        string        `json:"slug"`
	Status      domain.Status `json:"status"`
	Description string        `json:"description"`
	Content     string        `json:"content"`
	PublishedAt time.Time     `json:"published_at"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"udpated_at"`
}

type PostUpdate struct {
	ID          string        `param:"id" validate:"required"`
	Status      domain.Status `json:"status" validate:"omitempty,oneof=draft published archived"`
	Title       string        `json:"title" validate:"omitempty,required,min=5,max=128"`
	Author      string        `json:"author" validate:"omitempty,required,min=3,max=64"`
	Slug        string        `json:"slug" validate:"omitempty,required"`
	Description string        `json:"description" validate:"omitempty,required,min=5,max=255"`
	Content     string        `json:"content" validate:"omitempty,required,min=10"`
}

type PostsOut struct {
	Data []*PostOut `json:"data"`
}

type GetPostParams struct {
	ID string `param:"id" validate:"required"`
}
