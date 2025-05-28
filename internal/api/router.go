package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (api *Api) RegisterRoutes() {
	api.Echo.POST("/signup", api.SignUpHandler)
	api.Echo.POST("/signin", api.SignInHandler)
	api.Echo.POST("/logout", api.AuthMiddleware(api.LogoutHandler))
	api.Echo.GET("/protected", api.AuthMiddleware(api.protectedHandler))
}

func (api *Api) protectedHandler(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "protected route accessed"})
}
