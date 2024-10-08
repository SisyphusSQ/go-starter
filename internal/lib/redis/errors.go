package redis

import "github.com/redis/go-redis/v9"

var (
	// Nil record not found error
	Nil = redis.Nil
	// ErrNotFound record not found error
	ErrNotFound = redis.Nil
)
