package libs

import (
	"go.uber.org/fx"

	gormv2 "go-starter/internal/lib/gorm"
	"go-starter/internal/lib/redis"
)

var GlobalModule = fx.Provide(
	//log.New,
	gormv2.New,
	redis.New,
)
