package v1

import "github.com/labstack/echo/v4"

func InitAppRouter(handler *echo.Echo) {
	handler.GET("/go-template-echo/v1/test", getTest)
}
