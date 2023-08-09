package auth

import "github.com/labstack/echo/v4"

func (h *handler) Route(g *echo.Group) {
	g.GET("/status", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "User Service Oke"})
	})
	g.POST("/signup-users", h.RegisterUsers)
	g.POST("/request-otp", h.RequestOtp)
	g.POST("/verify-otp", h.VerifyOtp)
	g.POST("/loginbypin", h.LoginPin)
	g.POST("/admin-login", h.LoginAdmin)
	g.POST("/confirm-reset", h.ConfirmReset)
	g.POST("/reset-session", h.ResetDevice)
}
