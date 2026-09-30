package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"github.com/SisyphusSQ/go-starter/v2/config"
	"github.com/SisyphusSQ/go-starter/v2/internal/controller"
	"github.com/SisyphusSQ/go-starter/v2/internal/cron"
	"github.com/SisyphusSQ/go-starter/v2/internal/health"
	"github.com/SisyphusSQ/go-starter/v2/internal/http"
	libs "github.com/SisyphusSQ/go-starter/v2/internal/lib"
	"github.com/SisyphusSQ/go-starter/v2/internal/lib/log"
	"github.com/SisyphusSQ/go-starter/v2/internal/repository"
	"github.com/SisyphusSQ/go-starter/v2/internal/service"
	"github.com/SisyphusSQ/go-starter/v2/utils"
)

var configure string

var (
	httpCmd = &cobra.Command{
		Use:   "http",
		Short: "Start Http REST API",
		RunE:  initHTTP,
	}
)

func initHTTP(cmd *cobra.Command, _ []string) error {
	c, err := config.Load(configure)
	if err != nil {
		return err
	}
	if err = log.New(c); err != nil {
		return err
	}
	defer func() { _ = log.Sync() }()

	app := fx.New(inject(c))
	if err = app.Err(); err != nil {
		return fmt.Errorf("build application graph: %w", err)
	}

	runCtx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	startCtx, cancelStart := context.WithTimeout(runCtx, c.ContextTimeout)
	defer cancelStart()
	if err = app.Start(startCtx); err != nil {
		return fmt.Errorf("start application: %w", err)
	}
	var runtimeErr error
	select {
	case <-runCtx.Done():
	case signal := <-app.Wait():
		if signal.ExitCode != 0 {
			runtimeErr = fmt.Errorf("application stopped unexpectedly")
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), c.Server.ShutdownTimeout)
	defer cancel()
	if err = app.Stop(shutdownCtx); err != nil {
		return fmt.Errorf("stop application: %w", err)
	}
	return runtimeErr
}

func inject(c config.Config) fx.Option {
	return fx.Options(
		fx.Supply(c),
		fx.Provide(utils.NewTimeoutContext),
		fx.Provide(health.New),
		libs.Module(c),
		repository.Module,
		service.Module(c),
		cron.Module,
		controller.Module(c),
		fx.Invoke(func(*http.Server) {}),
		http.Module,
	)
}
