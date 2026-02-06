package resp

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"go-starter/utils"
	"go-starter/vars"
)

type Resp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func SuccessResp(data any) *Resp {
	return &Resp{
		Code:    vars.SUCCESS,
		Message: vars.GetMsg(vars.SUCCESS),
		Data:    data,
	}
}

func ErrorResp(err error) *Resp {
	return &Resp{
		Code:    vars.InternalERROR,
		Message: err.Error(),
		Data:    nil,
	}
}

func AssertErrResp(err string) *Resp {
	return &Resp{
		Code:    vars.InternalERROR,
		Message: err,
		Data:    nil,
	}
}

func CustomResp(code int, msg string, data any) *Resp {
	return &Resp{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}

type ResponseError struct {
	Message string `json:"message"`
}

type SimpleResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func CommErrResp(c echo.Context, err error) error {
	return c.JSON(utils.GetStatusCode(err), ErrorResp(err))
}

func CommSuccResp(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, SuccessResp(data))
}
