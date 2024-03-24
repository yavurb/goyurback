package ui

import "time"

type ProjectIn struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Tags         []string `json:"tags"`
	ThumbnailURL string   `json:"thumbnail_url"`
	WebsiteURL   string   `json:"website_url"`
	Live         bool     `json:"live"`
	PostId       int32    `json:"post_id"`
}

type ProjectOut struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Tags         []string  `json:"tags"`
	ThumbnailURL string    `json:"thumbnail_url"`
	WebsiteURL   string    `json:"website_url"`
	Live         bool      `json:"live"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	// PostId       string     `json:"post_id"` // TODO: Should be the post's public id
}

type GetProjectParam struct {
	ID string `param:"id" validate:"required"`
}

type ProjectsOut struct {
	Data []*ProjectOut `json:"data"`
}
