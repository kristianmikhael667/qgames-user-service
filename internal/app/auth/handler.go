package auth

import (
	"main/internal/dto"
	"main/internal/factory"
	"main/package/util/response"

	"github.com/labstack/echo/v4"
)

type handler struct {
	service Service
}

func NewHandler(f *factory.Factory) *handler {
	return &handler{
		service: NewService(f),
	}
}

func (h *handler) RegisterUsers(c echo.Context) error {
	payload := new(dto.RegisterUsersRequestBody)
	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	users, err := h.service.RegisterUsers(c.Request().Context(), payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	return response.SuccessResponse(users).Send(c)
}
