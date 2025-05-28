package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/SisyphusSQ/golib/models/vo/base_vo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go-starter/internal/lib/log"
	"go-starter/vars"
)

type EchoMiddleware struct {
}

func (e *EchoMiddleware) CORS(h echo.HandlerFunc) echo.HandlerFunc {
	cors := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.DELETE, echo.GET, echo.POST, echo.PUT, echo.OPTIONS, echo.HEAD, echo.PATCH},
	})
	return cors(h)
}

func (e *EchoMiddleware) Recover(h echo.HandlerFunc) echo.HandlerFunc {
	r := middleware.Recover()
	return r(h)
}

func (e *EchoMiddleware) Logger(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if strings.Contains(c.Request().RequestURI, "swagger") {
			return h(c)
		}
		log.Logger.Info("Enter method: [%s], uri: [%s], userAgent: [%s]", c.Request().Method, c.Request().RequestURI, c.Request().UserAgent())
		return h(c)
	}
}

func (e *EchoMiddleware) JWT(hf echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		uri := c.Request().RequestURI
		if strings.Compare(uri, "/") == 0 || strings.Compare(uri, "/health") == 0 ||
			strings.Contains(uri, "/swagger") {
			return hf(c)
		}
		jwtStr := c.Request().Header.Get("Authorization")
		auths := strings.Split(jwtStr, " ")
		if strings.ToUpper(auths[0]) != "BEARER" || auths[1] == "" {
			return c.JSON(http.StatusUnauthorized, base_vo.AssertErrResp("认证失败"))
		}
		// todo check jwt token
		// todo set jwt info in echo context
		return hf(c)
	}
}

func (e *EchoMiddleware) AccessAuth(hf echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		uri := c.Request().RequestURI
		if strings.Compare(uri, "/") == 0 || strings.Compare(uri, "/health") == 0 ||
			strings.Contains(uri, "swagger") {
			log.Logger.Debug("Directly enter to controller")
			return hf(c)
		}

		accessKey := c.Request().Header.Get("access_key")
		secretKey := c.Request().Header.Get("secret_key")
		if accessKey != vars.AccessKey && secretKey != vars.SecretKey {
			return c.JSON(http.StatusUnauthorized, base_vo.AssertErrResp("认证失败"))
		}

		return hf(c)
	}
}

func (e *EchoMiddleware) ErrorHandler(err error, c echo.Context) {
	var report *echo.HTTPError
	ok := errors.As(err, &report)
	if !ok {
		report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	log.Logger.Info("Leave method: [%s], uri: [%s], userAgent: [%s], got err: %v", c.Request().Method, c.Request().RequestURI, c.Request().UserAgent(), report.Message)
	c.Echo().DefaultHTTPErrorHandler(err, c)
}

func InitMiddleware() *EchoMiddleware {
	return &EchoMiddleware{}
}
