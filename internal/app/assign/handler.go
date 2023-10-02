package assign

import (
	"main/internal/dto"
	"main/internal/factory"
	"main/internal/pkg/util"
	"main/package/util/response"
	"net/http"

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

func (h *handler) EditAssign(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	_, err := util.ParseJWTToken(authHeader)

	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Unauthorized, err).Send(c)
	}

	payload := new(dto.ReqAssign)
	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.EditAssign(c, c.Request().Context(), payload)
	if err != nil && result == false {
		return response.CustomErrorBuilder(400, "Error", "Failed Edit Assign").Send(c)

	}
	return response.CustomSuccessBuilder(http.StatusOK, "Success", "Edit Assign Success", nil).Send(c)

}
