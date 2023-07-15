package role

import (
	dto "main/internal/dto/users_req_res"
	"main/internal/middleware"
	"main/internal/pkg/util"

	"github.com/labstack/echo/v4"
)

func (h *handler) Route(g *echo.Group) {
	g.Use(middleware.JWTMiddleware(dto.JWTClaims{}, util.JWT_SECRET))
	g.POST("/create-role", h.CreateRole)
}
