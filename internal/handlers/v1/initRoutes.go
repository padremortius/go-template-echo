package v1

import (
	"github.com/padremortius/go-template-echo/internal/config"
	"github.com/padremortius/go-template-echo/internal/storage"
	"github.com/padremortius/go-template-echo/pkgs/svclogger"

	"github.com/labstack/echo/v4"
)

type (
	v1Routes struct {
		cfg   config.Config
		log   svclogger.Log
		store storage.Storage
	}
)

func InitAppRouter(handler *echo.Echo, aCfg config.Config, aLog svclogger.Log, aStore storage.Storage) {
	v1 := v1Routes{cfg: aCfg, log: aLog, store: aStore}
	handler.GET("/go-template-echo/v1/test", v1.getTest)
}
