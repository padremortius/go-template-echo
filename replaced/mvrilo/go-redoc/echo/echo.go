package echoredoc

import (
	"github.com/labstack/echo/v5"
	"github.com/mvrilo/go-redoc"
)

func New(doc redoc.Redoc) echo.MiddlewareFunc {
	handle := doc.Handler()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			handle(ctx.Response(), ctx.Request())
			if resp, err := echo.UnwrapResponse(ctx.Response()); err == nil {
				if resp.Committed {
					return nil
				}
			}
			return next(ctx)
		}
	}
}
