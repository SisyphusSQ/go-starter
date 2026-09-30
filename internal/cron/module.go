package cron

import (
	"go.uber.org/fx"
)

var Module = fx.Invoke(
	NewCron,
)
