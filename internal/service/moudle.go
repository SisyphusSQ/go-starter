package service

import "go.uber.org/fx"

var Module = fx.Provide(
	NewAuditClusterService,
)
