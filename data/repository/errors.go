package repository

import "github.com/pkg/errors"

var (
	ErrCacheMissed   = errors.New("cache missed")
	ErrTypeAssertion = errors.New("type assertion failed")
)
