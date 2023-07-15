package http

import (
	"main/internal/app/users/auth"
	"main/internal/app/users/role"
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
	v1 := e.Group("/user-service")
	auth.NewHandler(f).Route(v1.Group("/users"))
	role.NewHandler(f).Route(v1.Group("/roles"))

	// v2 := e.Group("/product-service")

	// v3 := e.Group("/transaction-service")

	// v4 := e.Group("/wallet-service")

	// v5 := e.Group("/payment-service")

}
