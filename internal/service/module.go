package service

import (
	"github.com/SisyphusSQ/go-starter/v2/config"
	"go.uber.org/fx"

	"github.com/SisyphusSQ/go-starter/v2/internal/service/common_srv"
)

func Module(cfg config.Config) fx.Option {
	options := []fx.Option{fx.Provide()}

	if cfg.Lark.Enabled {
		options = append(options, fx.Provide(common_srv.NewLarkService), fx.Invoke(func(common_srv.LarkService) {}))
	}

	if cfg.Prometheus.Enabled {
		options = append(options, fx.Provide(common_srv.NewPrometheusService), fx.Invoke(func(common_srv.PrometheusService) {}))
	}

	return fx.Options(options...)
}
