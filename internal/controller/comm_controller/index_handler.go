package comm_controller

import (
	"net/http"

	"github.com/SisyphusSQ/golib/models/vo/base_vo"
	"github.com/labstack/echo/v4"

	"go-starter/utils/timeutil"
)

type IndexController struct{}

func InitIndexController(e *echo.Echo) {
	controller := &IndexController{}

	e.GET("/health", controller.Health)
	e.GET("/", controller.Health)
}

func (i *IndexController) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, base_vo.SuccessResp(timeutil.CSTLayoutString()))
}
