package user

import (
	dto "main/internal/dto"
	"main/internal/factory"
	"main/internal/pkg/util"
	pkgdto "main/package/dto"
	"main/package/util/response"
	res "main/package/util/response"
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

func (h *handler) Get(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	_, err := util.ParseJWTToken(authHeader)
	if err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.Unauthorized, err).Send(c)
	}

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

	return res.CustomSuccessBuilder(http.StatusOK, result.Data, "Get employees success", &result.PaginationInfo).Send(c)
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

	result, sc, msg, err := h.service.GetUserDetail(c.Request().Context(), token.Roles, token.Uuid)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.CustomSuccessBuilder(int(sc), result.Data, msg, nil).Send(c)
}
