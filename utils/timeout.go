package utils

import (
	"time"

	"github.com/spf13/viper"
)

func NewTimeoutContext() time.Duration {
	timeout := time.Duration(viper.GetInt("contextTimeout")) * time.Second

	return timeout
}
