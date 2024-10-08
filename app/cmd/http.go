package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"go-starter/config"
	"go-starter/internal/controller"
	_ "go-starter/internal/cron"
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

	if c.Key.SecretKey != "" {
		vars.SecretKey = c.Key.SecretKey
	}

	if c.Key.AccessKey != "" {
		vars.AccessKey = c.Key.AccessKey
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
		//cron.Module,
		controller.Module,
		http.Module,
	)
}
