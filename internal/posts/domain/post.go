package domain

import (
	"reflect"
	"time"
)

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

func (p Post) Compare(p2 Post) bool {
	p.CreatedAt = p2.CreatedAt
	p.UpdatedAt = p2.UpdatedAt
	p.PublishedAt = p2.PublishedAt

	return reflect.DeepEqual(p, p2)
}

type PostCreate struct {
	PublicID    string
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
