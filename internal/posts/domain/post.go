package domain

import "time"

type Status string

const (
	Draft     Status = "draft"
	Published Status = "published"
	Archived  Status = "archived"
)

type Post struct {
	ID          int
	Title       string
	Author      string
	Slug        string
	Status      Status
	Description string
	Content     string
	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PostCreate struct {
	Title       string
	Author      string
	Slug        string
	Description string
	Content     string
}

type PostUpdate struct {
	Title       string
	Author      string
	Slug        string
	Description string
	Content     string
	Status      Status
}
