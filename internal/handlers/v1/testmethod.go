package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// @Summary Test method
// @Description Test method
// @Produce json
// @Success 200 {object} JSONResult
// @Router /go-template-echo/v1/test [get]
// @Tags v1
func (v1 *v1Routes) getTest(c echo.Context) error {
	return c.JSON(http.StatusOK, &JSONResult{Code: http.StatusOK, Message: "Test complete!"})
}
