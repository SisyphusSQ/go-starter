package controller

import (
	"github.com/SisyphusSQ/go-starter/v2/config"
	"go.uber.org/fx"
)

func Module(cfg config.Config) fx.Option {
	var options []fx.Option

	return fx.Options(options...)
}
