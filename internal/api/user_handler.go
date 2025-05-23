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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := api.UserUseCase.CreateUser(c.Request().Context(), userReq); err != nil {
		apiErr := errors.ToAPIError(err)
		return c.JSON(apiErr.StatusCode, map[string]string{"error": apiErr.Message})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "user created"})
}
