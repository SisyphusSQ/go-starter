package cron

import (
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewCron,
)
