package libs

import (
	"github.com/SisyphusSQ/go-starter/v2/config"
	"go.uber.org/fx"

	gormv2 "github.com/SisyphusSQ/go-starter/v2/internal/lib/gorm"

	"github.com/SisyphusSQ/go-starter/v2/internal/lib/mongodb"

	"github.com/SisyphusSQ/go-starter/v2/internal/lib/redis"
)

// Module 仅装配显式启用的外部组件。
func Module(cfg config.Config) fx.Option {
	var options []fx.Option

	if cfg.Database.Enabled {
		options = append(options, fx.Provide(gormv2.New, fx.Annotate(gormv2.Readiness, fx.ResultTags(`group:"readiness"`))))
	}

	if cfg.MongoDB.Enabled {
		options = append(options, fx.Provide(mongodb.New, fx.Annotate(mongodb.Readiness, fx.ResultTags(`group:"readiness"`))))
	}

	if cfg.Redis.Enabled {
		options = append(options, fx.Provide(redis.New, fx.Annotate(redis.Readiness, fx.ResultTags(`group:"readiness"`))))
	}

	return fx.Options(options...)
}
