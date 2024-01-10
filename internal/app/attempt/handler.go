package attempt

import (
	dto "main/internal/dto"
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

func (h *handler) ResetOtpHandler(c echo.Context) error {
	phoneNumber := new(dto.RequestReset)

	if err := c.Bind(&phoneNumber); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.NotFound, err).Send(c)
	}
	if err := c.Validate(phoneNumber); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	msg, sc, err := h.service.ResetOtpService(c.Request().Context(), phoneNumber)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	if sc != 200 && sc != 201 {
		return response.CustomErrorBuilder(sc, "Reset OTP Error", msg).Send(c)
	}

	return response.CustomSuccessBuilder(sc, "Success Reset OTP To Number "+phoneNumber.Phone, msg, nil).Send(c)
}

func (h *handler) ResetPinHandler(c echo.Context) error {
	phoneNumber := new(dto.RequestReset)

	if err := c.Bind(&phoneNumber); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.NotFound, err).Send(c)
	}
	if err := c.Validate(phoneNumber); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	msg, sc, err := h.service.ResetPinService(c.Request().Context(), phoneNumber)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	if sc != 200 && sc != 201 {
		return response.CustomErrorBuilder(sc, "Reset PIN Error", msg).Send(c)
	}

	return response.CustomSuccessBuilder(sc, "Success Reset PIN To Number "+phoneNumber.Phone, msg, nil).Send(c)
}
