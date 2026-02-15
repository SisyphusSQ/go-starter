package http

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/SisyphusSQ/golib/models/vo/base_vo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go-starter/config"
	"go-starter/internal/lib/log"
	redisv9 "go-starter/internal/lib/redis"
	"go-starter/utils"
	"go-starter/vars"
)

type EchoMiddleware struct {
	config config.Config
	cache  *redisv9.Client
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
		log.Logger.Info("Enter method: [%s], uri: [%s], userAgent: [%s]", c.Request().Method, c.Request().RequestURI, c.Request().UserAgent())
		return h(c)
	}
}

func (e *EchoMiddleware) JWT(hf echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		uri := c.Request().URL.Path
		if e.isPublicURI(uri) {
			return hf(c)
		}

		token, err := e.extractBearerToken(c.Request().Header.Get("Authorization"))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, base_vo.AssertErrResp("认证失败"))
		}

		claims, err := utils.ParseToken(token, e.config.Key.JWT.Secret)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, base_vo.AssertErrResp("认证失败"))
		}

		if e.cache == nil {
			return c.JSON(http.StatusUnauthorized, base_vo.AssertErrResp("认证失败"))
		}

		cacheKey := fmt.Sprintf("jwt:user:%d", claims.UserID)
		cacheToken, err := e.cache.Get(c.Request().Context(), cacheKey).Result()
		if err != nil {
			if errors.Is(err, redisv9.Nil) {
				return c.JSON(http.StatusUnauthorized, base_vo.AssertErrResp("认证失败"))
			}
			return c.JSON(http.StatusUnauthorized, base_vo.AssertErrResp("认证失败"))
		}
		if cacheToken != token {
			return c.JSON(http.StatusUnauthorized, base_vo.AssertErrResp("认证失败"))
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		return hf(c)
	}
}

func (e *EchoMiddleware) AccessAuth(hf echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		uri := c.Request().RequestURI
		if strings.Compare(uri, "/") == 0 || strings.Compare(uri, "/health") == 0 {
			log.Logger.Debug("Directly enter to controller")
			return hf(c)
		}

		accessKey := c.Request().Header.Get("access_key")
		secretKey := c.Request().Header.Get("secret_key")
		if accessKey != vars.AccessKey || secretKey != vars.SecretKey {
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

func (e *EchoMiddleware) isPublicURI(uri string) bool {
	return uri == "/" ||
		uri == "/health" ||
		uri == "/login" ||
		strings.Contains(uri, "/swagger")
}

func (e *EchoMiddleware) extractBearerToken(authorization string) (string, error) {
	auths := strings.SplitN(authorization, " ", 2)
	if len(auths) != 2 {
		return "", errors.New("invalid authorization header")
	}
	if !strings.EqualFold(auths[0], "Bearer") {
		return "", errors.New("invalid authorization type")
	}
	if strings.TrimSpace(auths[1]) == "" {
		return "", errors.New("empty token")
	}
	return auths[1], nil
}

func InitMiddleware(config config.Config, cache *redisv9.Client) *EchoMiddleware {
	return &EchoMiddleware{
		config: config,
		cache:  cache,
	}
}
