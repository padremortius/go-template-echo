module github.com/mvrilo/go-redoc/echo

go 1.25.0

require (
	github.com/labstack/echo/v5 v5.4.4
	github.com/mvrilo/go-redoc v0.1.5-0.20240119235626-7cdf2a730a86
)

require github.com/goccy/go-json v0.10.5 // indirect

replace github.com/mvrilo/go-redoc => ../
