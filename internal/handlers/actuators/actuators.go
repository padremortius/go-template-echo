package actuators

import (
	"go-template-echo/internal/config"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mvrilo/go-redoc"
	echoredoc "github.com/mvrilo/go-redoc/echo"
)

type (
	Health struct {
		Status string
	}

	BaseRoutes struct {
		cfg config.Config
	}
)

func (b *BaseRoutes) getHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, &Health{Status: "up"})
}

func (b *BaseRoutes) getInfo(c echo.Context) error {
	return c.JSON(http.StatusOK, &b.cfg.Version)
}

func (b *BaseRoutes) getEnv(c echo.Context) error {
	return c.JSON(http.StatusOK, &b.cfg)
}

func InitBaseRouter(handler *echo.Echo, aCfg config.Config) {
	bRoutes := BaseRoutes{cfg: aCfg}

	// K8s probe
	handler.GET("/health", bRoutes.getHealth)

	//info about service
	handler.GET("/info", bRoutes.getInfo)

	//env
	handler.GET("/env", bRoutes.getEnv)

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
