package domain

import (
	"reflect"
	"time"
)

type Project struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	PublicID     string
	Name         string
	Description  string
	ThumbnailURL string
	WebsiteURL   string
	Tags         []string
	ID           int32
	PostID       int32
	Live         bool
}

func (p Project) Compare(p2 Project) bool {
	p.CreatedAt = p2.CreatedAt
	p.UpdatedAt = p2.UpdatedAt

	return reflect.DeepEqual(p, p2)
}

type ProjectCreate struct {
	PublicID     string
	Name         string
	Description  string
	ThumbnailURL string
	WebsiteURL   string
	Tags         []string
	PostID       int32
	Live         bool
}
