package api

import (
	"github.com/josevitorrodriguess/gochat/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Api struct {
	Echo        *echo.Echo
	UserUseCase usecase.UserUseCase
}

func NewApi(userUc usecase.UserUseCase) *Api {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	return &Api{
		Echo:        e,
		UserUseCase: userUc,
	}
}

func (api *Api) Start(port string) error {
	return api.Echo.Start(port)
}
