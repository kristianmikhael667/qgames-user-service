package attempt

import (
	dto "main/internal/dto"
	"main/internal/middleware"
	"main/internal/pkg/util"

	"github.com/labstack/echo/v4"
)

func (h *handler) Route(g *echo.Group) {
	g.GET("/status", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "USER SERVICE OK"})
	})

	g.Use(middleware.JWTMiddleware(dto.JWTClaims{}, util.JWT_SECRET))
	g.POST("/reset-otp", h.ResetOtpHandler)
	g.POST("/reset-pin", h.ResetPinHandler)
}
