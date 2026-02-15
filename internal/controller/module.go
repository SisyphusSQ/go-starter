package controller

import (
	"go.uber.org/fx"

	"go-starter/internal/controller/comm_controller"
	"go-starter/internal/controller/example_controller"
)

var Module = fx.Invoke(
	comm_controller.InitIndexController,
	example_controller.InitUserController,
	example_controller.InitUserMongoController,
)
