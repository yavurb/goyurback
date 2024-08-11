package domain

import "time"

type Status string

const (
	Draft     Status = "draft"
	Published Status = "published"
	Archived  Status = "archived"
)

type Post struct {
	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PublicID    string
	Title       string
	Author      string
	Slug        string
	Status      Status
	Description string
	Content     string
	ID          int32
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
