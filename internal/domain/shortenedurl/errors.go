package shortenedurl

import "errors"

var (
	ErrNotFound  = errors.New("shortened url not found")
	ErrDuplicate = errors.New("shortened url already exists")
)
