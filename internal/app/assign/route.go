package assign

import (
	"main/internal/dto"
	"main/internal/middleware"
	"main/internal/pkg/util"

	"github.com/labstack/echo/v4"
)

func (h *handler) Route(g *echo.Group) {
	g.Use(middleware.JWTMiddleware(dto.JWTClaims{}, util.JWT_SECRET))
	g.PUT("/assigns-roles", h.EditAssign)
}
