package http

import (
	"main/internal/app/auth"
	"main/internal/app/permission"
	"main/internal/app/role"
	"main/internal/app/user"
	"main/internal/factory"
	"main/package/util"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

func NewHttp(e *echo.Echo, f *factory.Factory) {
	e.Validator = &util.CustomValidator{Validator: validator.New()}

	e.GET("/status", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "OKE"})
	})
	v1 := e.Group("/")
	auth.NewHandler(f).Route(v1.Group("user-service"))
	role.NewHandler(f).Route(v1.Group("user-service"))
	user.NewHandler(f).Route(v1.Group("user-service"))
	permission.NewHandler(f).Route(v1.Group("user-service"))

	// v2 := e.Group("/product-service")

	// v3 := e.Group("/transaction-service")

	// v4 := e.Group("/wallet-service")

	// v5 := e.Group("/payment-service")

}
