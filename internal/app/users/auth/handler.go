package auth

import (
	dto "main/internal/dto/users_req_res"
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

// func (h *handler) LoginUser(c echo.Context)

func (h *handler) RequestOtp(c echo.Context) error {
	phoneNumber := new(dto.CheckPhoneReqBody)
	if err := c.Bind(phoneNumber); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.NotFound, err).Send(c)
	}
	if err := c.Validate(phoneNumber); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	checkphone, status, err := h.service.RequestOtp(c.Request().Context(), phoneNumber)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	if status == true {
		return response.CustomSuccessBuilder(201, checkphone, "Waiting OTP Send", nil).Send(c)
	} else {
		return response.CustomSuccessBuilder(200, checkphone, "Waiting OTP Send", nil).Send(c)
	}

}
