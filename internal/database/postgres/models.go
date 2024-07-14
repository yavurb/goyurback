// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package postgres

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusArchived  PostStatus = "archived"
)

func (e *PostStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PostStatus(s)
	case string:
		*e = PostStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for PostStatus: %T", src)
	}
	return nil
}

type NullPostStatus struct {
	PostStatus PostStatus
	Valid      bool // Valid is true if PostStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPostStatus) Scan(value interface{}) error {
	if value == nil {
		ns.PostStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PostStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPostStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PostStatus), nil
}

type Apikey struct {
	ID        int32
	PublicID  string
	Name      string
	Key       string
	Revoked   bool
	RevokedAt pgtype.Timestamp
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type Post struct {
	ID          int32
	PublicID    string
	Title       string
	Author      string
	Content     string
	Description string
	Slug        string
	Status      PostStatus
	PublishedAt pgtype.Timestamp
	CreatedAt   pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
}

type Project struct {
	ID           int32
	PublicID     string
	Name         string
	Description  string
	Tags         []string
	ThumbnailUrl string
	WebsiteUrl   string
	Live         bool
	CreatedAt    pgtype.Timestamp
	UpdatedAt    pgtype.Timestamp
	PostID       pgtype.Int4
}
