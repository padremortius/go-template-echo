package actuators

import (
	"go-template-echo/internal/config"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mvrilo/go-redoc"
	echoredoc "github.com/mvrilo/go-redoc/echo"
)

type (
	Health struct {
		Status string
	}
)

func getHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, &Health{Status: "up"})
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

	//redoc
	doc := redoc.Redoc{
		SpecFile: "spec/docs.json",
		SpecPath: "/docs.json",
		DocsPath: "/redoc",
		Options: map[string]any{
			"disableSearch": true,
			"theme": map[string]any{
				"colors":     map[string]any{"primary": map[string]any{"main": "#297b21"}},
				"typography": map[string]any{"headings": map[string]any{"fontWeight": "600"}},
				"sidebar":    map[string]any{"backgroundColor": "lightblue"},
			},
			"decorator": map[string]any{},
		},
	}
	handler.Use(echoredoc.New(doc))
}
