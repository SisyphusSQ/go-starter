package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"go-starter/config"
	"go-starter/internal/controller"
	"go-starter/internal/http"
	libs "go-starter/internal/lib"
	"go-starter/internal/repository"
	"go-starter/internal/service"
	"go-starter/utils"
)

var (
	httpCmd = &cobra.Command{
		Use:   "http",
		Short: "Start Http REST API",
		Run:   initHTTP,
	}
)

func initHTTP(cmd *cobra.Command, args []string) {
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
		controller.Module,
		http.Module,
	)
}
