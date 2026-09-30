package utils

import (
	"time"

	"github.com/SisyphusSQ/go-starter/v2/config"
)

func NewTimeoutContext(c config.Config) time.Duration {
	return c.ContextTimeout
}
