package middleware

import (
	dto "main/internal/dto/users_req_res"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func LogMiddlewares(e *echo.Echo) {
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           `[${time_rfc3339}] ${status} ${method} ${host}${uri} ${latency_human}` + "\n",
		CustomTimeFormat: "2006/01/02 15:04:05",
	}))
}

func JWTMiddleware(claims dto.JWTClaims, signingKey []byte) echo.MiddlewareFunc {
	config := middleware.JWTConfig{
		Claims:     &dto.JWTClaims{},
		SigningKey: signingKey,
	}
	return middleware.JWTWithConfig(config)
}
