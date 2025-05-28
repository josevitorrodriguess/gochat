package api

import (
	"encoding/gob"
	"net/http"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/josevitorrodriguess/gochat/internal/storage/redis"
	"github.com/josevitorrodriguess/gochat/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Api struct {
	Echo        *echo.Echo
	UserUseCase usecase.UserUseCase

	RedisStorage *redis.RedisStorage
	Sessions     *scs.SessionManager
}

func NewApi(userUc usecase.UserUseCase) *Api {
	gob.Register(uuid.UUID{})
	e := echo.New()

	redisStore, err := redis.NewDefault()
	if err != nil {
		e.Logger.Fatal("could not connect to Redis:", err)
	}

	s := scs.New()
	s.Store = redisstore.New(redisStore.GetPool())
	s.Lifetime = 24 * time.Hour
	s.Cookie.Name = "gochat_session"
	s.Cookie.HttpOnly = true
	s.Cookie.Secure = false
	s.Cookie.SameSite = http.SameSiteLaxMode

	e.Use(echo.WrapMiddleware(s.LoadAndSave))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	return &Api{
		Echo:        e,
		UserUseCase: userUc,
		Sessions:    s,
	}
}

func (api *Api) Start(port string) error {
	return api.Echo.Start(port)
}
