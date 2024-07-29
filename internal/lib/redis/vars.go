package redis

import (
	"github.com/go-redis/redis/v8"
)

type Client = redis.Client
type Cmder = redis.Cmder
type Cmdable = redis.Cmdable
type ScanIterator = redis.ScanIterator
type Pipeline = redis.Pipeline
type PubSub = redis.PubSub
type Pipeliner = redis.Pipeliner
