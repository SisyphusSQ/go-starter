package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"go-starter/internal/models/resp"
	"go-starter/utils/timeutil"
)

type IndexController struct {
}

func InitIndexController(e *echo.Echo) {
	controller := &IndexController{}
	e.GET("/health", controller.Health)
	e.GET("/", controller.Health)
}

func (index *IndexController) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, resp.SimpleResponse{
		Message: "success",
		Data:    timeutil.CSTLayoutString(),
	})
}
