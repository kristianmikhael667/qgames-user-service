package auth

import (
	"fmt"
	dto "main/internal/dto"
	"main/internal/factory"
	"main/internal/pkg/util"
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

func (h *handler) RequestOtp(c echo.Context) error {
	phoneNumber := new(dto.CheckPhoneReqBody)

	if err := c.Bind(&phoneNumber); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.NotFound, err).Send(c)
	}
	if err := c.Validate(phoneNumber); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	checkphone, sc, _, err := h.service.RequestOtp(c.Request().Context(), phoneNumber)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	if sc == 403 {
		return response.CustomErrorBuilder(sc, "Device Login", checkphone).Send(c)
	}

	if sc == 201 {
		return response.CustomSuccessBuilder(sc, checkphone, "New User", nil).Send(c)
	} else if sc == 200 {
		return response.CustomSuccessBuilder(sc, checkphone, "Old users", nil).Send(c)
	} else {
		return response.CustomErrorBuilder(sc, checkphone, "Error").Send(c)
	}
}

func (h *handler) VerifyOtp(c echo.Context) error {
	bodyVerify := new(dto.RequestPhoneOtp)
	if err := c.Bind(bodyVerify); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.NotFound, err).Send(c)
	}
	if err := c.Validate(bodyVerify); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	checkverify, msg, sc, _ := h.service.VerifyOtp(c.Request().Context(), bodyVerify)
	if sc != 201 {
		return response.CustomErrorBuilder(int(sc), "error", msg).Send(c)
	}

	return response.CustomSuccessBuilder(int(sc), checkverify, msg, nil).Send(c)
}

func (h *handler) LoginPin(c echo.Context) error {
	payloads := new(dto.LoginByPin)
	if err := c.Bind(payloads); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.NotFound, err).Send(c)
	}
	if err := c.Validate(payloads); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, msg, sc, err := h.service.LoginPin(c.Request().Context(), payloads)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	if sc != 201 {
		return response.CustomErrorBuilder(sc, "Error", msg).Send(c)
	} else {
		return response.CustomSuccessBuilder(sc, result, msg, nil).Send(c)
	}
}

func (h *handler) CheckPin(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	token, err := util.ParseJWTToken(authHeader)
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Unauthorized, err).Send(c)
	}

	payloads := new(dto.CheckPin)
	if err := c.Bind(payloads); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.NotFound, err).Send(c)
	}
	if err := c.Validate(payloads); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	isPin, sc, _, err := h.service.CheckPin(c.Request().Context(), token, payloads)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	if sc != 201 {
		return response.CustomErrorBuilder(sc, "error", "Wrong PIN").Send(c)
	} else {
		return response.CustomSuccessBuilder(sc, isPin, "True PIN", nil).Send(c)
	}
}

func (h *handler) LoginAdmin(c echo.Context) error {
	payloads := new(dto.LoginAdmin)
	if err := c.Bind(payloads); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.NotFound, err).Send(c)
	}
	if err := c.Validate(payloads); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, msg, sc, err := h.service.LoginAdmin(c.Request().Context(), payloads)
	if err != nil {
		return response.CustomErrorBuilder(sc, msg, msg).Send(c)
	}

	if sc != 201 {
		return response.CustomErrorBuilder(sc, msg, msg).Send(c)
	} else {
		return response.CustomSuccessBuilder(sc, result, msg, nil).Send(c)
	}
}

func (h *handler) ConfirmReset(c echo.Context) error {
	phoneNumber := new(dto.CheckSession)

	if err := c.Bind(&phoneNumber); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.NotFound, err).Send(c)
	}
	if err := c.Validate(phoneNumber); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	msg, sc, err := h.service.ConfirmReset(c.Request().Context(), phoneNumber)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	if sc != 201 {
		return response.CustomErrorBuilder(sc, "Error", msg).Send(c)
	}
	return response.CustomSuccessBuilder(sc, "reset-device", msg, nil).Send(c)
}

func (h *handler) ResetDevice(c echo.Context) error {
	session := new(dto.ReqSessionReset)

	if err := c.Bind(&session); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.NotFound, err).Send(c)
	}
	if err := c.Validate(session); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, msg, sc, err := h.service.ResetDevice(c.Request().Context(), session)
	fmt.Println(result)
	if err != nil {
		return response.CustomErrorBuilder(sc, "Error", msg).Send(c)
	}

	if sc != 201 && sc != 205 {
		return response.CustomErrorBuilder(sc, "Error", msg).Send(c)
	}
	return response.CustomSuccessBuilder(sc, result, msg, nil).Send(c)
}

func (h *handler) RefreshToken(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	token, err := util.ParseJWTToken(authHeader)
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Unauthorized, err).Send(c)
	}

	isToken, sc, _, err := h.service.RefreshToken(c.Request().Context(), token)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	if sc != 201 {
		return response.CustomErrorBuilder(sc, "error", "Error Refresh Token").Send(c)
	} else {
		return response.CustomSuccessBuilder(sc, isToken, "Success Refresh Token", nil).Send(c)
	}
}
