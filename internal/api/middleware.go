package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (api *Api) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		fmt.Printf("Cookies recebidos: %+v\n", c.Request().Cookies())

		userID := api.Sessions.GetString(c.Request().Context(), "user_id")

		if userID == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		c.Set("user_id", userID)
		return next(c)
	}
}
