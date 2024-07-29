package redis

import "github.com/go-redis/redis/v8"

var (
	// Nil record not found error
	Nil = redis.Nil
	// ErrNotFound record not found error
	ErrNotFound = redis.Nil
)
