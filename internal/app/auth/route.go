package auth

import (
	dto "main/internal/dto"
	"main/internal/middleware"
	"main/internal/pkg/util"

	"github.com/labstack/echo/v4"
)

func (h *handler) Route(g *echo.Group) {
	g.GET("/status", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "Tuhan aku nih jodoh yang body slim itu rine atau devina ?"})
	})
	g.POST("/signup-users", h.RegisterUsers)
	g.POST("/request-otp", h.RequestOtp)
	g.POST("/verify-otp", h.VerifyOtp)
	g.POST("/loginbypin", h.LoginPin)
	g.POST("/admin-login", h.LoginAdmin)
	g.POST("/confirm-reset", h.ConfirmReset)
	g.POST("/reset-session", h.ResetDevice)

	// Use Token
	g.Use(middleware.JWTMiddleware(dto.JWTClaims{}, util.JWT_SECRET))
	g.POST("/check-pin", h.CheckPin)
	g.POST("/refresh-token", h.RefreshToken)
}
