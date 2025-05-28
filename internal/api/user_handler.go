package api

import (
	"net/http"

	"github.com/josevitorrodriguess/gochat/internal/errors"
	"github.com/josevitorrodriguess/gochat/internal/models/request"
	"github.com/labstack/echo/v4"
)

func (api *Api) SignUpHandler(c echo.Context) error {
	var userReq request.UserRequest

	if err := c.Bind(&userReq); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "invalid request"})
	}

	if err := api.UserUseCase.CreateUser(c.Request().Context(), userReq); err != nil {
		apiErr := errors.ToAPIError(err)
		return c.JSON(apiErr.StatusCode, map[string]string{"error": apiErr.Message})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "user created"})
}

func (api *Api) SignInHandler(c echo.Context) error {
	var signInReq request.SignInRequest

	if err := c.Bind(&signInReq); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "invalid request"})
	}

	ok, err := signInReq.IsValid()
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	userID, err := api.UserUseCase.AuthenticateUser(c.Request().Context(), signInReq.Email, signInReq.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	err = api.Sessions.RenewToken(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to renew session token"})
	}

	api.Sessions.Put(c.Request().Context(), "user_id", userID.String())

	return c.JSON(http.StatusOK, map[string]string{"message": "user authenticated"})
}

func (api *Api) LogoutHandler(c echo.Context) error {
	err := api.Sessions.RenewToken(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to renew session token"})
	}

	api.Sessions.Remove(c.Request().Context(), "user_id")

	return c.JSON(http.StatusOK, map[string]string{"message": "user logged out"})
}
