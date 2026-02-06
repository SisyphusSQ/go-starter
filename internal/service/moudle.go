package service

import (
	"go.uber.org/fx"

	"go-starter/internal/service/common_srv"
	"go-starter/internal/service/example_srv"
)

var Module = fx.Provide(
	common_srv.NewLarkService,
	common_srv.NewPrometheusService,
	example_srv.NewAuditClusterService,
)
