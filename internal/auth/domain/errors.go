package domain

import "errors"

var ErrAPIKeyNotFound = errors.New("post not found")
var ErrAPIKeyInvalid = errors.New("invalid api key")
