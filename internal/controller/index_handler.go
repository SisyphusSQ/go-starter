package controller

import (
	"net/http"

	"github.com/SisyphusSQ/golib/models/vo/base_vo"
	"github.com/labstack/echo/v4"

	"go-starter/internal/cron"
	"go-starter/utils/timeutil"
)

type IndexController struct {
	cronSrv cron.Service
}

func InitIndexController(e *echo.Echo, cronSrv cron.Service) {
	controller := &IndexController{
		cronSrv: cronSrv,
	}

	e.GET("/health", controller.Health)
	e.GET("/", controller.Health)
	e.GET("/host", controller.Host)
}

func (i *IndexController) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, base_vo.SuccessResp(timeutil.CSTLayoutString()))
}

func (i *IndexController) Host(c echo.Context) error {
	return base_vo.CommSuccResp(c, i.cronSrv.IP())
}
