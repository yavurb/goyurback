package domain

import "errors"

var (
	ErrChikitoNotFound       = errors.New("no chikito was found")
	ErrPublicIDAlreadyExists = errors.New("public id already exists")
)
