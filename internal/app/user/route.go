package user

import (
	dto "main/internal/dto"
	"main/internal/middleware"
	"main/internal/pkg/util"

	"github.com/labstack/echo/v4"
)

func (h *handler) Route(g *echo.Group) {
	g.PUT("/register-user", h.UpdateUser)
	g.Use(middleware.JWTMiddleware(dto.JWTClaims{}, util.JWT_SECRET))
	g.GET("/users", h.Get)
	g.GET("/user/:id", h.GetUserId)
	g.GET("/myaccount", h.MyAccount)
	g.POST("/reset-pin", h.ResetPin)
	g.POST("/logout", h.Logout)
}
