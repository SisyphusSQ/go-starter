package example_controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/SisyphusSQ/golib/models/vo/base_vo"
	"github.com/labstack/echo/v4"

	"go-starter/internal/models/vo"
	"go-starter/internal/service/example_srv"
	"go-starter/utils"
)

type UserController struct {
	userService example_srv.UserService
}

func InitUserController(e *echo.Echo, userService example_srv.UserService) {
	controller := &UserController{
		userService: userService,
	}

	e.POST("/login", controller.Login)
	e.POST("/logout", controller.Logout)

	g := e.Group("/mysql/users")
	g.GET("/:id", controller.GetByID)
	g.GET("", controller.List)
	g.POST("", controller.Create)
	g.PUT("/:id", controller.Update)
	g.DELETE("/:id", controller.Delete)
}

func (u *UserController) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return base_vo.CommErrResp(c, utils.ErrBadParamInput)
	}

	user, err := u.userService.GetByID(c.Request().Context(), id)
	if err != nil {
		return base_vo.CommErrResp(c, err)
	}
	return base_vo.CommSuccResp(c, user)
}

func (u *UserController) List(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		return base_vo.CommErrResp(c, utils.ErrBadParamInput)
	}
	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil {
		return base_vo.CommErrResp(c, utils.ErrBadParamInput)
	}

	resp, err := u.userService.List(c.Request().Context(), page, pageSize)
	if err != nil {
		return base_vo.CommErrResp(c, err)
	}

	return base_vo.CommSuccResp(c, resp)
}

func (u *UserController) Create(c echo.Context) error {
	var req vo.CreateUserReq
	if err := c.Bind(&req); err != nil {
		return base_vo.CommErrResp(c, err)
	}

	user, err := u.userService.Create(c.Request().Context(), req)
	if err != nil {
		return base_vo.CommErrResp(c, err)
	}
	return base_vo.CommSuccResp(c, user)
}

func (u *UserController) Login(c echo.Context) error {
	var req vo.LoginReq
	if err := c.Bind(&req); err != nil {
		return base_vo.CommErrResp(c, err)
	}

	resp, err := u.userService.Login(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, example_srv.ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, base_vo.AssertErrResp("认证失败"))
		}
		return base_vo.CommErrResp(c, err)
	}
	return base_vo.CommSuccResp(c, resp)
}

func (u *UserController) Logout(c echo.Context) error {
	userID, ok := getUserID(c.Get("user_id"))
	if !ok {
		return c.JSON(http.StatusUnauthorized, base_vo.AssertErrResp("认证失败"))
	}

	resp, err := u.userService.Logout(c.Request().Context(), userID)
	if err != nil {
		return base_vo.CommErrResp(c, err)
	}
	return base_vo.CommSuccResp(c, resp)
}

func (u *UserController) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return base_vo.CommErrResp(c, utils.ErrBadParamInput)
	}

	var req vo.UpdateUserReq
	if err = c.Bind(&req); err != nil {
		return base_vo.CommErrResp(c, err)
	}

	resp, err := u.userService.Update(c.Request().Context(), id, req)
	if err != nil {
		return base_vo.CommErrResp(c, err)
	}
	return base_vo.CommSuccResp(c, resp)
}

func (u *UserController) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return base_vo.CommErrResp(c, utils.ErrBadParamInput)
	}

	resp, err := u.userService.Delete(c.Request().Context(), id)
	if err != nil {
		return base_vo.CommErrResp(c, err)
	}
	return base_vo.CommSuccResp(c, resp)
}

func getUserID(v any) (int64, bool) {
	switch value := v.(type) {
	case int64:
		return value, value > 0
	case int:
		id := int64(value)
		return id, id > 0
	case int32:
		id := int64(value)
		return id, id > 0
	case float64:
		id := int64(value)
		return id, id > 0
	case string:
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return 0, false
		}
		return id, id > 0
	default:
		return 0, false
	}
}
