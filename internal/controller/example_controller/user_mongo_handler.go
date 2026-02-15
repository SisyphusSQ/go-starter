package example_controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/SisyphusSQ/golib/models/vo/base_vo"
	"github.com/labstack/echo/v4"

	"go-starter/internal/lib/log"
	"go-starter/internal/models/vo"
	"go-starter/internal/repository/mongo/example_repo"
	"go-starter/internal/service/example_srv"
	"go-starter/utils"
)

type UserMongoController struct {
	userMongoService example_srv.UserMongoService
}

func InitUserMongoController(e *echo.Echo, userMongoService example_srv.UserMongoService) {
	if !userMongoService.IsAvailable() {
		if log.Logger != nil {
			log.Logger.Warnf("mongo user service unavailable, skip /mongo/users routes")
		}
		return
	}

	controller := &UserMongoController{
		userMongoService: userMongoService,
	}

	g := e.Group("/mongo/users")
	g.GET("/:id", controller.GetByID)
	g.GET("", controller.List)
	g.POST("", controller.Create)
	g.PUT("/:id", controller.Update)
	g.DELETE("/:id", controller.Delete)
}

func (u *UserMongoController) handleErr(c echo.Context, err error) error {
	if errors.Is(err, example_repo.ErrMongoUnavailable) {
		return c.JSON(http.StatusServiceUnavailable, base_vo.AssertErrResp("mongo service unavailable"))
	}
	return base_vo.CommErrResp(c, err)
}

func (u *UserMongoController) GetByID(c echo.Context) error {
	id := c.Param("id")
	user, err := u.userMongoService.GetByID(c.Request().Context(), id)
	if err != nil {
		return u.handleErr(c, err)
	}
	return base_vo.CommSuccResp(c, user)
}

func (u *UserMongoController) List(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		return base_vo.CommErrResp(c, utils.ErrBadParamInput)
	}
	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil {
		return base_vo.CommErrResp(c, utils.ErrBadParamInput)
	}

	resp, err := u.userMongoService.List(c.Request().Context(), page, pageSize)
	if err != nil {
		return u.handleErr(c, err)
	}

	return base_vo.CommSuccResp(c, resp)
}

func (u *UserMongoController) Create(c echo.Context) error {
	var req vo.CreateUserReq
	if err := c.Bind(&req); err != nil {
		return base_vo.CommErrResp(c, err)
	}

	user, err := u.userMongoService.Create(c.Request().Context(), req)
	if err != nil {
		return u.handleErr(c, err)
	}
	return base_vo.CommSuccResp(c, user)
}

func (u *UserMongoController) Update(c echo.Context) error {
	id := c.Param("id")

	var req vo.UpdateUserReq
	if err := c.Bind(&req); err != nil {
		return base_vo.CommErrResp(c, err)
	}

	resp, err := u.userMongoService.Update(c.Request().Context(), id, req)
	if err != nil {
		return u.handleErr(c, err)
	}
	return base_vo.CommSuccResp(c, resp)
}

func (u *UserMongoController) Delete(c echo.Context) error {
	id := c.Param("id")
	resp, err := u.userMongoService.Delete(c.Request().Context(), id)
	if err != nil {
		return u.handleErr(c, err)
	}
	return base_vo.CommSuccResp(c, resp)
}
