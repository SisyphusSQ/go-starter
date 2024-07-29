package http

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/fx"

	"go-starter/config"
)

var Module = fx.Provide(NewServer)

func NewServer(lifecycle fx.Lifecycle, config config.Config) *echo.Echo {
	instance := echo.New()
	middleware := InitMiddleware()

	instance.Use(middleware.CORS)
	instance.Use(middleware.Logger)
	instance.Use(middleware.Recover)
	instance.Use(middleware.AccessAuth)

	instance.HTTPErrorHandler = middleware.ErrorHandler

	instance.GET("/swagger/*", echoSwagger.WrapHandler)

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Start Http Server.")
			go func() {
				err := instance.Start(config.Server.Address)
				if err != nil {
					_ = fmt.Errorf("start Http Server error: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Stopping Http Server.")
			return instance.Shutdown(ctx)
		},
	})
	return instance
}
