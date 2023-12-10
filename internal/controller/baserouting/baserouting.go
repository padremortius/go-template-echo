package baserouting

import (
	"fmt"
	"go-template-echo/docs"
	"go-template-echo/internal/config"
	"go-template-echo/internal/controller/structs"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func getHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, &structs.Health{Status: "up"})
}

func getInfo(c echo.Context) error {
	return c.JSON(http.StatusOK, &config.Cfg.Version)
}

func getEnv(c echo.Context) error {
	return c.JSON(http.StatusOK, &config.Cfg)
}

func InitBaseRouter(handler *echo.Echo) {
	handler.Use(middleware.Recover())

	// K8s probe
	handler.GET("/health", getHealth)

	//info about service
	handler.GET("/info", getInfo)

	//env
	handler.GET("/env", getEnv)

	//swagger
	docs.SwaggerInfo.Version = config.Cfg.Version.Version
	docs.SwaggerInfo.BasePath = fmt.Sprint("/", config.Cfg.Name)
	handler.GET(fmt.Sprint(docs.SwaggerInfo.BasePath, "/swagger/*"), echoSwagger.WrapHandler)
}
