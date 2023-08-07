package auth

import (
	"fmt"
	dto "main/internal/dto"
	"main/internal/factory"
	"main/internal/pkg/util"
	"main/package/util/response"
	"time"

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

	fmt.Println("user id ", phoneNumber.Phone)
	fmt.Println("dev ", phoneNumber.DeviceId)

	checkphone, sc, status, err := h.service.RequestOtp(c.Request().Context(), phoneNumber)
	if err != nil {
		fmt.Println("sc ? ", sc, " ", checkphone)
		return response.ErrorResponse(err).Send(c)
	}

	if sc == 403 {
		return response.CustomErrorBuilder(sc, "Device Login", checkphone).Send(c)
	}

	if sc == 201 && status == true {
		return response.CustomSuccessBuilder(sc, checkphone, "New User", nil).Send(c)
	} else if sc == 200 && status == false {
		return response.CustomSuccessBuilder(sc, checkphone, "Old users", nil).Send(c)
	} else {
		return response.CustomSuccessBuilder(sc, checkphone, "Old users", nil).Send(c)
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
		return response.CustomErrorBuilder(int(sc), err.Error(), msg).Send(c)
	} else {
		return response.CustomSuccessBuilder(int(sc), result, msg, nil).Send(c)
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

func (h *handler) RevokeToken(c echo.Context) error {
	var blacklistMap = make(map[string]time.Time)
	authHeader := c.Request().Header.Get("Authorization")
	token, _ := util.ParseJWTToken(authHeader)

	blacklistMap[token.ExpiresAt.String()] = time.Now()

	return response.CustomSuccessBuilder(201, "Success Logout", "Revoke Token", nil).Send(c)

}
