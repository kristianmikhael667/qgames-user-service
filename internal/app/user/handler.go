package user

import (
	dto "main/internal/dto"
	"main/internal/factory"
	"main/internal/pkg/util"
	pkgdto "main/package/dto"
	"main/package/util/response"
	res "main/package/util/response"
	"net/http"
	"strings"

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

func (h *handler) Get(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	token, err := util.ParseJWTToken(authHeader)

	if err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.Unauthorized, err).Send(c)
	}

	if token.Admin {
		payload := new(pkgdto.SearchGetRequest)
		if err := c.Bind(payload); err != nil {
			return res.ErrorBuilder(&res.ErrorConstant.BadRequest, err).Send(c)
		}

		if err := c.Validate(payload); err != nil {
			return res.ErrorBuilder(&res.ErrorConstant.Validation, err).Send(c)
		}

		result, err := h.service.Find(c.Request().Context(), payload)
		if err != nil {
			return res.ErrorResponse(err).Send(c)
		}

		return res.CustomSuccessBuilder(http.StatusOK, result.Data, "Get users success", &result.PaginationInfo).Send(c)

	}
	return res.CustomErrorBuilder(403, "Forbidden", "Token not allowed to access").Send(c)
}

func (h *handler) GetUserId(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	_, err := util.ParseJWTToken(authHeader)
	if err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.Unauthorized, err).Send(c)
	}

	payload := new(pkgdto.ByIDRequest)

	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.FindIdUser(c.Request().Context(), payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	return response.CustomSuccessBuilder(200, result, "Get User ID", nil).Send(c)
}

func (h *handler) UpdateUser(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	token, err := util.ParseJWTToken(authHeader)
	if err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.Unauthorized, err).Send(c)
	}

	payload := new(dto.UpdateUsersReqBody)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	uid := token.Uuid
	users, sc, msg, err := h.service.UpdateUsers(c.Request().Context(), &pkgdto.ByUuidUsersRequest{Uid: uid}, payload)

	if err != nil || sc != 201 {
		return response.CustomErrorBuilder(int(sc), msg, err.Error()).Send(c)
	}

	return response.CustomSuccessBuilder(int(sc), users, msg, nil).Send(c)
}

func (h *handler) MyAccount(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	token, err := util.ParseJWTToken(authHeader)
	if err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.Unauthorized, err).Send(c)
	}

	for _, permission := range token.Permissions {
		hasCommonUser := strings.Contains(permission, "common-user")
		if hasCommonUser {
			result, sc, msg, err := h.service.GetUserDetail(c, c.Request().Context(), token.Roles, token.Uuid)
			if err != nil {
				return res.ErrorResponse(err).Send(c)
			}

			if sc == 403 || sc == 201 {
				return res.CustomErrorBuilder(sc, msg, "error").Send(c)
			}

			return res.CustomSuccessBuilder(int(sc), result, msg, nil).Send(c)
		}
	}
	return res.CustomErrorBuilder(403, "Forbidden", "Token not allowed to access").Send(c)
}

func (h *handler) ResetPin(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	token, err := util.ParseJWTToken(authHeader)
	if err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.Unauthorized, err).Send(c)
	}

	payload := new(dto.ConfirmPin)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	uid := token.Uuid
	users, sc, msg, err := h.service.ResetPin(c.Request().Context(), uid, token.Roles, payload)

	if sc != 201 {
		return response.CustomErrorBuilder(sc, msg, "Error").Send(c)
	}

	return response.CustomSuccessBuilder(sc, users, msg, nil).Send(c)
}

func (h *handler) Logout(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	token, err := util.ParseJWTToken(authHeader)
	if err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.Unauthorized, err).Send(c)
	}

	uid := token.Uuid
	msg, sc, err := h.service.Logout(c, c.Request().Context(), uid)

	if sc != 201 && sc != 200 {
		return response.CustomErrorBuilder(sc, msg, "Error").Send(c)
	}

	return response.CustomSuccessBuilder(sc, msg, "Success Logout", nil).Send(c)
}
