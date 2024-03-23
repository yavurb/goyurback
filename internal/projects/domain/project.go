package domain

import "time"

type Project struct {
	ID           int32
	PublicID     string
	Name         string
	Description  string
	Tags         []string
	ThumbnailURL string
	WebsiteURL   string
	Live         bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	PostID       int32
}

type ProjectCreate struct {
	Name         string
	Description  string
	Tags         []string
	ThumbnailURL string
	WebsiteURL   string
	Live         bool
	PostID       int32
}
