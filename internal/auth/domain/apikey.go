package domain

import "time"

type APIKey struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	KeyName  string
	Key      string
	PublicID string
	ID       int32
	Revoked  bool
}

type APIKeyCreate struct {
	KeyName string
	Key     string
}
