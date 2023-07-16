package permission

import "github.com/labstack/echo/v4"

func (h *handler) Route(g *echo.Group) {
	// g.Use(middleware.JWTMiddleware(dto.JWTClaims{}, util.JWT_SECRET))
	g.POST("/create-permission", h.CreatePermission)
}
