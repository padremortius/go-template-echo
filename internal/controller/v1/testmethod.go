package v1

import (
	"go-template-echo/internal/controller/structs"
	"net/http"

	"github.com/labstack/echo/v4"
)

// @Summary Test method
// @Description Test method
// @Produce json
// @Success 200 {object} structs.JSONResult
// @Router /v1/test [get]
// @Tags v1
func getTest(c echo.Context) error {
	return c.JSON(http.StatusOK, &structs.JSONResult{Code: http.StatusOK, Message: "Test complete!"})
}
