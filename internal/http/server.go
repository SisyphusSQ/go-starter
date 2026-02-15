package http

import (
	"context"
	"fmt"

	prom "github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	mid "github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"

	"go-starter/config"
	redisv9 "go-starter/internal/lib/redis"
	"go-starter/vars"
)

var Module = fx.Provide(NewServer)

func NewServer(lifecycle fx.Lifecycle, config config.Config, cache *redisv9.Client) *echo.Echo {
	instance := echo.New()
	middleware := InitMiddleware(config, cache)

	instance.Use(middleware.CORS)
	instance.Use(middleware.Logger)
	instance.Use(middleware.Recover)
	instance.Use(prom.NewMiddleware("go_starter"))

	switch config.Key.Type {
	case "basic":
		instance.Use(mid.BasicAuth(func(user string, password string, c echo.Context) (bool, error) {
			if user == vars.User && password == vars.Password {
				return true, nil
			}

			return false, nil
		}))
	case "key":
		instance.Use(middleware.AccessAuth)
	case "jwt":
		instance.Use(middleware.JWT)
	}

	instance.HTTPErrorHandler = middleware.ErrorHandler

	instance.GET("/metrics", prom.NewHandler())

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
