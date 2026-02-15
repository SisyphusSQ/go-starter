package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"go-starter/config"
	"go-starter/internal/controller"
	"go-starter/internal/cron"
	"go-starter/internal/http"
	libs "go-starter/internal/lib"
	"go-starter/internal/lib/log"
	"go-starter/internal/repository"
	"go-starter/internal/service"
	"go-starter/utils"
	"go-starter/vars"
)

var configuare string

var (
	httpCmd = &cobra.Command{
		Use:   "http",
		Short: "Start Http REST API",
		Run:   initHTTP,
	}
)

func initHTTP(cmd *cobra.Command, args []string) {
	config.SetConfigFile(configuare)
	config.InitConfig()
	c := config.NewConfig()
	log.New(c)
	defer log.Logger.Sync()

	switch c.Key.Type {
	case "basic":
		if c.Key.Basic.User != "" && c.Key.Basic.Password != "" {
			vars.User = c.Key.Basic.User
			vars.Password = c.Key.Basic.Password
		} else {
			panic("basic auth required")
		}
	case "key":
		if c.Key.AK.SecretKey != "" && c.Key.AK.AccessKey != "" {
			vars.SecretKey = c.Key.AK.SecretKey
			vars.AccessKey = c.Key.AK.AccessKey
		} else {
			panic("key auth required")
		}
	case "jwt":
		if c.Key.JWT.Secret == "" || c.Key.JWT.Expire <= 0 {
			panic("jwt config required")
		}
	default:
		panic("auth required")
	}

	fx.New(inject()).Run()
}

func inject() fx.Option {
	return fx.Options(
		fx.Provide(
			config.NewConfig,
			utils.NewTimeoutContext,
		),
		libs.GlobalModule,
		repository.Module,
		service.Module,
		cron.Module,
		controller.Module,
		http.Module,
	)
}
